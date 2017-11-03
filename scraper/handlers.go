/*
    handlers just turn the logic of slicing up HTML into little bits™
*/
package scraper

import (
    "strings"
    "bytes"
    "errors"
    "strconv"
    "time"

    "github.com/anaskhan96/soup"
)

//makes schedule object from string of Dates and Times, Musn't include TBA
//NOTE: Doesn't set location
func makeSchedule(dateTimePtr string) *Schedule {

		//treat all the newlines and carriage returns
		sched := &Schedule{}
		dateTime := strings.Trim(dateTimePtr, "\n\r ")
		var dtl int = len(dateTime)
		tba := strings.HasPrefix(dateTime, "TBA")

		//treat TBA strings differently
		if !tba && dtl > 0 {

				// startTime and endTime will always come at the same index
				startT, _ := time.Parse("03:04 pm", dateTime[dtl - 20: dtl - 12])
				endT, _ := time.Parse("03:04 pm", dateTime[dtl - 8: dtl])
				sched.StartTime = startT.Format("15:04")
				sched.EndTime = endT.Format("15:04")
				sched.Days = dateTime[:dtl - 26]
		} else {
        sched.StartTime, sched.EndTime, sched.Days = "TBA", "TBA", "TBA"
    }

		return sched
}

//makes schedule, but also returns previous location
func makeSchedulePrev(schedPtr string) (string, *Schedule) {

		// since the line starts with the previous schedule's data, this is a bit tricky
		// finding the index of the start of the new schedule from the back
		splitSched := strings.Split(strings.Replace(schedPtr, "\n", " ", -1), " ")
		var index int
		for i, str := range splitSched {
				if str == "from" {
					index = strings.Index(schedPtr, splitSched[i - 1])
				}
		}

		//means the new schedule's TBA, but we can return the last one
		if index == 0 {

				//if there is no "from", then they're both TBA and there's no reason to add the second
			 	return strings.Trim(schedPtr, "\n "), nil

		} else {

				// since now we know where schedule 1 loc is
				prev := schedPtr[:index - 1]

				// splits the schedule string apart from sched location
				secondSched := strings.Trim(schedPtr[index:], "\n\r ")
				sched := makeSchedule(secondSched)

				// adds it to list of schedules
				return prev, sched
		}
}

// Writes the data from the first part to the Section struct
func handleHead(root soup.Root) (*Section, error) {

		// starting from the tr
		a := root.Find("a");

		if a.Error != nil {
				return nil, errors.New("Can't parse empty element.")
		}

		// a.text would be something like ECE 9940 - 001
		// split data into [DeptName, CourseNum, -, Section]
		courseInfo := strings.Split(a.Text(), " ")

		// move along to dirtier data:
		// Art & Culture Exp Learn Com CRN: 33914 Enrollment: FULL 17 students.
		// or 									 								  Enrollment: 0 of 30 students.
		// I'm gonna go ahead and work around this API...
		el := a.Pointer.Parent.LastChild

		//reset a
		a = soup.Root{el, el.Data, nil}
		rawMeta := el.Data
		if len(rawMeta) < 5 {
				return nil, errors.New("Can't parse an empty string")
		}

		//extract all that precious metadata
		name, crn, enrolled, size := extractMetadata(rawMeta)

		//setup section, make attributes too
		section := &Section{
				CRN: crn,
				Name: name,
				Section: courseInfo[3],
				Enrolled: enrolled,
				Size: size,
				Class: courseInfo[1],
				Department: courseInfo[0],
		}

		return section, nil
}

//extract name, crn, enrollment data from string
func extractMetadata(str string) (name, crn string, enrolled, size int) {

		i := strings.Index(str, "CRN")

		//get name and cut it out (accounting for whitespace)
		name = str[:i - 1]

		//extract CRN which should be
		crn = str[i + 5:i + 10]

		//extracts last section like "Enrollment: FULL 17 students. "
		enrollment := str[i + 12:]
		enrollmentSplit := strings.Split(enrollment, " ")
		var num int

		//parse int out of enrollment string
		num, _ = strconv.Atoi(enrollmentSplit[len(enrollmentSplit) - 3])
		size = num

		//set enrolled last because it depends on size
		if enrollmentSplit[2] == "FULL" {
				enrolled = num
		} else {
				enrolled, _ = strconv.Atoi(enrollmentSplit[1])
		}

		return
}

// The body section is a little easier to parse, it's labeled by classnames
//start, end, instructor string, restrictions, attributes []string
func handleCourseBody(root soup.Root) (sched []*Schedule, rest, cmt string, instr []*Instructor, attr, reqs[]string) {

		//if we're starting on an error, hop out - usually no sibling error
		if root.Error != nil {
				return
		}

		// these basically mark all important data here
		pointers := root.Find("td").FindAll("span")

		// root is the wrong element sometimes, causing an error here
		// guessing that fieldlabeltext isn't exclusive to sections
		if len(pointers) > 0 {

				// loop through all the sections
				for i := 0; i < len(pointers); i++ {

						// start with dateTime string
						data := pointers[i].Text()

						// will write data into the Section
						switch strings.Trim(data, " ") {
						case "Days:":
								sched = handleDateLoc(pointers[i])
						case "Comment:":
								cmt = handleComment(pointers[i])
						case "Instructors:":
								instr = handleInstructor(pointers[i])
						case "Prerequisites:":
								reqs = handlePrereqs(pointers[i])
						case "Restrictions:":
								rest = handleRestrictions(pointers[i])
						case "Attributes:":
								attr = handleAttributes(pointers[i])
						}
				}
		}

		return
}

// Should only work where classes have up to 2 locations & meeting times
// "Days:" will be followed by a text node describing the date & then Location
func handleDateLoc(root soup.Root) []*Schedule {
		schedules := make([]*Schedule, 0)
		dateTimePtr := root.FindNextSibling()
		dateTime := strings.Trim(dateTimePtr.Pointer.Data, "\n\r ")

		// works for the head, we just tac it on
		sched := makeSchedule(dateTime)

		// adds schedule to schedules
		schedules = append(schedules, sched)

		//Go up, get all children of element <b>, identifying num of classes
		parent := &soup.Root{
				Pointer: root.Pointer.Parent,
				NodeValue: root.Pointer.Parent.Data,
				Error: nil,
		}

		// if we're dealing with a section with 2 schedules
		for i, ptr := range parent.FindAll("b") {
				schedData := ptr.FindNextSibling().Pointer.Data
				prev, sch := makeSchedulePrev(schedData)
				schedules[i].Location = emptyToTBA(prev)
				if sch != nil {

						// adds schedule to schedules
						schedules = append(schedules, sched)
				} else {

						//if we have a TBA, the others will be TBA
						break
				}
		}
		return schedules
}

// "Comment:" will occasionally have a description for a class
func handleComment(root soup.Root) string {
		data := root.FindNextSibling().Pointer.Data
		cleaned := strings.Replace(strings.Replace(data, "\n", " ", -1), "&nbsp;", " ", -1)
		return strings.Trim(cleaned, " ")
}

// "Instructors:" will be followed by a list of teachers
func handleInstructor(root soup.Root) []*Instructor {
		teachers := make([]*Instructor, 0)

		//instructors are best found in <a> tags, so let's step up to parent
		parent := &soup.Root{
				Pointer: root.Pointer.Parent,
				NodeValue: root.Pointer.Parent.Data,
				Error: nil,
		}

		teacherPtrs := parent.FindAll("a")
		for _, teacher := range teacherPtrs {
				t := &Instructor{
						Name: teacher.Attrs()["target"],
						Id: teacher.Attrs()["href"], 		 //teacher emails are unique to them
				}
        if t.Name == "" {
            continue
        }
				teachers = append(teachers, t)
		}

		return teachers
}

// "Restrictions" will scrape out generic restriction data (usually just UG/G)
func handleRestrictions(root soup.Root) string {
		root = root.FindNextSibling()
		var buffer bytes.Buffer
		for root.Error == nil && strings.Index(root.Pointer.Data, "University Alliance") == -1 {
				if root.Pointer.Data != "br" {
  					cleanStr := strings.Replace(root.Pointer.Data, "&nbsp;", " ", -1)
  					buffer.WriteString(cleanStr)
  					buffer.WriteString("\n")
				}
				root = root.FindNextSibling()
		}
    //if we've written to the buffer, we've seen "University Alliance" -- just not written it
    if buffer.Len() > 0 {
        buffer.WriteString("University Alliance")
    }
		return buffer.String()
}

// "Prerequisites", kinda important
func handlePrereqs(root soup.Root) []string {
		prereqs := make([]string, 0)

		//instructors are best found in <a> tags, so let's step up to parent
		parent := &soup.Root{
				Pointer: root.Pointer.Parent,
				NodeValue: root.Pointer.Parent.Data,
				Error: nil,
		}

		prereqPtrs := parent.FindAll("a")
		for _, prereq := range prereqPtrs {
				data := prereq.Text()
				if data != "" {
						prereqs = append(prereqs, data)
			 	}
		}

		return prereqs
}

// "Attributes:" sections are extremely important
func handleAttributes(root soup.Root) []string {
		data := root.FindNextSibling().Pointer.Data
		cleaned := strings.Replace(strings.Replace(data, "\n", " ", -1), "&nbsp;", " ", -1)
		split := strings.Split(strings.Trim(cleaned, " "), " ")
		return split
}

func emptyToTBA(locTime string) string {
    if locTime == "" {
      return "TBA"
    }
    return locTime
}
