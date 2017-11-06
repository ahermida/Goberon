/*
    Package scraper is a utility for extracting data from the course data html page

		Structure - <tr> come in pairs as sibling elements
		<table class="datadisplaytable">
			<tbody>
				<tr></tr> // Metadata like CRN, Course name, and Section
				<tr></tr>	// Instructor restrictions, Prereqs, Attributes etc...
				<tr></tr> // Metadata again
				<tr></tr> // All that other data
				...
			</tbody>
		</table>
*/
package scraper

import (
		"fmt"
		"io/ioutil"

		"github.com/anaskhan96/soup"
		"github.com/ahermida/Goberon/config"
)

var (
		Registrar *CourseData
)

type CourseData struct {
		Instructors map[string]*Instructor
		Courses 		map[string]*Course
		Departments map[string]*Department
		CI 					map[string]*Section
}

// Go ahead and index the registry of data
func Index() (int, error) {
		fmt.Println("Loading resources")
		err := InitDB()
		if err != nil {
			return 0, err
		}
		coursesHTML, _ := ioutil.ReadFile(config.Local.DefaultFN)
		doc := soup.HTMLParse(string(coursesHTML))
		tables := doc.FindAll("table", "class", "datadisplaytable")
		registrar := start(tables)
		if config.CatSecret != nil {
				cat := handleCatalog()
				indexCatalog(cat, registrar)
		}
		Commit(registrar)
		return 0, nil
}

func Commit(reg *CourseData) {

		//Index all data into empty DB
		WriteSections(reg.CI)
		WriteCourses(reg.Courses)
		WriteInstructors(reg.Instructors)
		WriteDepartments(reg.Departments)
}


// Every department has a caption tag followed by a data table body - so we start here
func start(tables []soup.Root) *CourseData {

		// keep the current active registrar exported
		Registrar = &CourseData{
				CI: 				  map[string]*Section{},
				Instructors:  map[string]*Instructor{},
				Courses: 		  map[string]*Course{},
				Departments:  map[string]*Department{},
		}

		total := 0

		// go over tables, extracting data as we can
		for _, table := range tables {

				// for each table found
				// <table>
				//       \
				//         -> <tbody>
				//                  \
				//                    -> <tr> <tr> <tr> ...
				if tr := table.Find("tr"); tr.Error == nil {

						//writes course data into registrar passed in
						total += handleBody(tr, Registrar.CI)
				}
		}

		//By now, weve seen all the courses so let's start building other tables
	  for _, v := range Registrar.CI {

				//build instructors from sections
				for _, t := range v.Instructors {

						//if instructor doesn't exist, make em
						if Registrar.Instructors[t.Id] == nil {
								t.Sections = make([]string, 0)
								Registrar.Instructors[t.Id] = t
						}
						v.Instructor = append(v.Instructor, t.Name)
						Registrar.Instructors[t.Id].Sections = append(Registrar.Instructors[t.Id].Sections, v.CRN)
				}

				//build department if not exists
				if Registrar.Departments[v.Department] == nil {
						Registrar.Departments[v.Department] = &Department{
								Name: v.Department,
								Courses: make([]string, 0),
						}
				}

				//get course Id
				cid := v.Department + v.Class

				//build course or add to it
				if Registrar.Courses[cid] == nil {
						Registrar.Courses[cid] = &Course{
								ID: cid,
								Name: v.Name,
								Class: v.Class,
								Department: v.Department,
								Sections: []string{v.CRN},
						}

						//add course to Department
						Registrar.Departments[v.Department].Courses = append(Registrar.Departments[v.Department].Courses, cid)
				} else {
						Registrar.Courses[cid].Sections = append(Registrar.Courses[cid].Sections, v.CRN)
				}
		}

		return Registrar
}

// Handles the <tr>s within <tbody>, must come in pairs of two
func handleBody(root soup.Root, ci map[string]*Section) int {
		var i = 0

		// keep chugging along until there's no more to parse
		for root.Error == nil {

				// mostly populated section
				section, _ := handleHead(root)

				// don't move over unless we actually can
				root = root.FindNextElementSibling()

				//being called 2k+ times
				sched, rest, cmt, instr, attr, reqs := handleCourseBody(root)

				if section != nil {
						section.Schedule = sched
						section.Comment = cmt
						section.Instructors = instr
						section.Prereqs = reqs
						section.Attributes = attr
						section.Restrictions = rest

						//stores all sections in map
						ci[section.CRN] = section

						//we've seen another course!
						i++
				}

				// next element since we're doing this in pair
				if root.Error == nil {
						root = root.FindNextElementSibling()
				}
		}

		return i
}

// parse through the catalog, writing different sections to CourseData
func handleCatalog() []*CatalogData{
		catHTML, _ := ioutil.ReadFile(config.Local.CatFN)
		doc := soup.HTMLParse(string(catHTML))
		tables := doc.FindAll("table", "class", "datadisplaytable")
		li := make([]*CatalogData, 0)
		for _, table := range tables {
				if tr := table.Find("tr"); tr.Error == nil {
						for _, c := range handleCatBody(tr) {
								li = append(li, c)
						}
				}
		}

		//give em what they NEED
		return li
}

// write course data to catalog -- Catalog gives us more than are offered any given semester
func indexCatalog(cat []*CatalogData, cd *CourseData) {
		for _, c := range cat {
				if cd.Courses[c.Course] != nil {

						//apply to course
						cd.Courses[c.Course].Description = c.Description
						cd.Courses[c.Course].Credits = c.Credits
						for _, s := range cd.Courses[c.Course].Sections {
								cd.CI[s].Description = c.Description
								cd.CI[s].Credits = c.Credits
						}
				}
		}
}
