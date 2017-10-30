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
		Instructors map[string]*Instructor
		Courses map[string]*Course
		Departments map[string]*Department
		CI map[string]*Section
)

func init() {
		Instructors = map[string]*Instructor{}
	 	Courses	= map[string]*Course{}
		Departments = map[string]*Department{}
		CI = map[string]*Section{}
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
		captions := doc.FindAll("caption")
		tables := doc.FindAll("table", "class", "datadisplaytable")
		start(captions, tables)
		Commit()
		return 0, nil
}

func Commit() {

		//Index all data into empty DB
		WriteSections(CI)
		WriteCourses(Courses)
		WriteInstructors(Instructors)
		WriteDepartments(Departments)
}


// Every department has a caption tag followed by a data table body - so we start here
func start(captions, tables []soup.Root) error {
		majors := map[string]bool{}

		// pass over map of captions to get list of majors without parsing data just yet
		for _, caption := range captions {
				if !majors[caption.Text()] {
						majors[caption.Text()] = true
				}
		}

		// update ourselves on Number of Majors
		fmt.Printf("Found %v departments", len(majors))
		var total = 0
		// over captions again, this time traversing each caption's neighbor
		for _, table := range tables {

				// for each table found
				// <table>
				//       \
				//         -> <tbody>
				//                  \
				//                    -> <tr> <tr> <tr> ...
				if tr := table.Find("tr"); tr.Error == nil {
						i := handleBody(tr, CI)
						total += i
				}
		}

		//By now, weve seen all the courses so let's start building other tables
	  for _, v := range CI {

				//build instructors from sections
				for _, t := range v.Instructors {
						v.Instructor = append(v.Instructor, t.Name)
						Instructors[t.Id].Sections = append(Instructors[t.Id].Sections, v.CRN)
				}

				//build department if not exists
				if Departments[v.Department] == nil {
						Departments[v.Department] = &Department{
								Name: v.Department,
								Courses: make([]string, 0),
						}
				}

				//get course Id
				cid := v.Department + v.Class

				//build course or add to it
				if Courses[cid] == nil {
						Courses[cid] = &Course{
								ID: cid,
								Name: v.Name,
								Class: v.Class,
								Department: v.Department,
								Sections: []string{v.CRN},
						}

						//add course to Department
						Departments[v.Department].Courses = append(Departments[v.Department].Courses, cid)
				} else {
						Courses[cid].Sections = append(Courses[cid].Sections, v.CRN)
				}
		}

		fmt.Printf("\nTotal Courses Indexed: %d\n", total)

		return nil
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
