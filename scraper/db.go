package scraper

import (
    "log"
    "fmt"

    r "gopkg.in/gorethink/gorethink.v3"
    "gopkg.in/cheggaaa/pb.v1"
    "github.com/ahermida/Goberon/config"
)

var Session *r.Session
var DB r.Term
var tables []string
var primaries []string

func init() {
    tables = []string{"sections", "teachers", "courses", "departments"}
    primaries = []string{"crn", "id", "id", "name"}
}

//not running this in init to prevent slow startup
func InitDB() error {
    var err error
    DB, Session, err = startDB()
    if err != nil {
      return err
    }
    err = dropTables(DB, Session)
    if err != nil {
      return err
    }
    err = createTables(DB, Session)
    return err
}

//CLI drop function
func Drop() error {
    db, session, err := startDB()
    if err != nil {
      return err
    }
    err = dropTables(db, session)
    return err
}

func startDB() (r.Term, *r.Session, error) {
    var session *r.Session
    var err error
    session, err = r.Connect(r.ConnectOpts{
        Address: config.Network.RdbAddress,
        Database: config.Network.RDB,
    })
    if err != nil {
        return DB, nil, err
    }

    //reference db
    DB = r.DB(config.Network.RDB)
    return DB, session, nil
}

func dropTables(db r.Term, session *r.Session) error {

    for _, t := range tables {
        DB.TableDrop(t).Run(session)
    }
    return nil
}

func createTables(db r.Term, session *r.Session) error {
    for i, t := range tables {
        opts := r.TableCreateOpts{
            PrimaryKey: primaries[i],
        }
        _, err := DB.TableCreate(t, opts).Run(session)
        if err != nil {
          return err
        }
    }
    return nil
}

// writes sections to Rethink
func WriteSections(ci map[string]*Section) error {
    fmt.Println("Writing Sections...")
    bar := pb.New(len(ci))
    bar.SetMaxWidth(80)
    bar.Start()
    for _, s := range ci {
        if s != nil {
            bar.Increment()
            _, err := DB.Table("sections").Insert(s).RunWrite(Session)
            if err != nil {
                log.Panic(err)
                return err
            }

        }
    }
    bar.Finish()
    return nil
}

// writes instructors to Rethink
func WriteInstructors(instructors map[string]*Instructor) error {
    fmt.Println("Writing Instructors...")
    bar := pb.New(len(instructors))
    bar.SetMaxWidth(80)
    bar.Start()
    for _, i := range instructors {
        bar.Increment()
        _, err := DB.Table("teachers").Insert(i).RunWrite(Session)
        if err != nil {
            return err
        }
    }
    bar.Finish()
    return nil
}

// writes courses to Rethink
func WriteCourses(courses map[string]*Course) error {
    fmt.Println("Writing Courses...")
    bar := pb.New(len(courses))
    bar.SetMaxWidth(80)
    bar.Start()
    for _, c := range courses {
        bar.Increment()
        _, err := DB.Table("courses").Insert(c).RunWrite(Session)
        if err != nil {
            return err
        }
    }
    bar.Finish()
    return nil
}

// writes Departments to Rethink
func WriteDepartments(departments map[string]*Department) error {
    fmt.Println("Writing Departments...")
    bar := pb.New(len(departments))
    bar.SetMaxWidth(80)
    bar.Start()
    for _, d := range departments {
        bar.Increment()
        _, err := DB.Table("departments").Insert(d).RunWrite(Session)
        if err != nil {
            return err
        }
    }
    bar.Finish()
    return nil
}
