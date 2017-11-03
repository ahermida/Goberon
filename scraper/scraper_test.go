package scraper

import (
    "testing"
)

const testHTML = `
<html>
  <head>
    <title>test</title>
  </head>
  <body>
    <div class="pagebodydiv">
      <table class="datadisplaytable">
        <caption class="captiontext">Sections FoundVillanova Experience -- Spring 2018</caption>
        <tbody>
          <tr>
            <th class="ddlabel" scope="row"><a href="/pls/bannerprd/bwckctlg.p_disp_course_detail?cat_term_in=201830&amp;subj_code_in=VEXP&amp;crse_numb_in=0001">VEXP 0001 - AC1 </a>Art &amp; Culture Exp Learn Com CRN: 33913 Enrollment: 15 of 16 students. </th>
            <td><a href="http://www.bkstr.com/webapp/wcs/stores/servlet/booklookServlet ?bookstore_id-1=1349&amp;term_id-1=201830&amp;crn-1=33913">Book List</a></td>
          </tr>
          <tr>
            <td class="dddefault">
            <br>
            <span class="fieldlabeltext">Days: </span>
            R
            from 02:30 pm to 03:45 pm
            <b> Location: </b>TBA
            <br>
            <span class="fieldlabeltext">Instructors: </span>
            Tom  DeMarco <a href="mailto:tom.demarco@villanova.edu" target="Tom DeMarco"><img src="/wtlgifs/web_email.gif" align="middle" alt="E-mail" class="headerImg" title="E-mail" name="web_email" hspace="0" vspace="0" border="0" height="28" width="28"></a> (<abbr title="Primary">P</abbr>)
            <br>
            <span class="fieldlabeltext">Restrictions:</span>
            Must be enrolled in one of the following Levels:&nbsp; &nbsp; &nbsp;
            <br>
            &nbsp; &nbsp; &nbsp; Undergraduate
            <br>
            May not be enrolled in one of the following Campuses:&nbsp; &nbsp; &nbsp;
            <br>
            &nbsp; &nbsp; &nbsp; University Alliance
            <br>
            <br>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </body>
</html>
`
func testMetadata() {

}
