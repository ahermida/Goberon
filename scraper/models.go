/*
    Models for rethinkDB
*/
package scraper

// might as well model out all these different tables
type Section struct {
    CRN           string         `gorethink:"crn"`           //33919
    Instructors   []*Instructor  `gorethink:"-"`             // -- Used for indexing
    Instructor    []string       `gorethink:"instructors"`   //[Daniel Joyce]
    Attributes    []string       `gorethink:"attributes"`    //Core Science
    Restrictions  string         `gorethink:"restrictions"`  //University Alliance
    Prereqs       []string       `gorethink:"prereqs"`       //CSC-1051
    Name          string         `gorethink:"name"`          //Algorithms & Data Structures
    Comment       string         `gorethink:"comment"`       //Learn all about Linked Lists!
    Section       string         `gorethink:"section"`       //001
    Class         string         `gorethink:"class"`         //1052
    Enrolled      int            `gorethink:"enrolled"`      //07
    Size          int            `gorethink:"size"`          //30
    Department    string         `gorethink:"department"`    //CSC
    Schedule      []*Schedule    `gorethink:"schedule"`      //[{10:30, 11:45, MWF, Mendel G86}]
}

// 1052 is taught by plenty of teachers
type Course struct {
    ID          string     `gorethink:"id,omitempty"` //Department+Class (should be unique)
    Department  string     `gorethink:"department"`
    Name        string     `gorethink:"name"`
    Class       string     `gorethink:"class"`
    Sections    []string   `gorethink:"sections"`

}

// instructor teaches many sections
type Instructor struct {
    Name      string     `gorethink:"name"`
    Id        string     `gorethink:"id"`
    Sections    []string   `gorethink:"sections"`
}

// Departments can have many teachers, but teachers don't always belong to a department
type Department struct {
    Name      string     `gorethink:"name"`
    Courses   []string   `gorethink:"courses"`
}

// days, times, and location of a class
type Schedule struct {
		Days       string  `gorethink:"days"`       //MWF
		StartTime  string  `gorethink:"startTime"`  //15:15
		EndTime    string  `gorethink:"endTime"`    //16:30
		Location   string  `gorethink:"location"`   //Mendel G86
}
