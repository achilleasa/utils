// Copyright 2013 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package series_test

import (
	"io/ioutil"
	"path/filepath"
	"sort"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/utils/os"
	"github.com/juju/utils/series"
)

func (s *supportedSeriesSuite) TestSeriesVersion(c *gc.C) {
	// There is no distro-info on Windows or CentOS.
	if os.HostOS() != os.Ubuntu {
		c.Skip("This test is only relevant on Ubuntu.")
	}
	vers, err := series.SeriesVersion("precise")
	if err != nil && err.Error() == `invalid series "precise"` {
		c.Fatalf(`Unable to lookup series "precise", you may need to: apt-get install distro-info`)
	}
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(vers, gc.Equals, "12.04")
}

func (s *supportedSeriesSuite) TestSupportedSeries(c *gc.C) {
	d := c.MkDir()
	filename := filepath.Join(d, "ubuntu.csv")
	err := ioutil.WriteFile(filename, []byte(distInfoData), 0644)
	c.Assert(err, jc.ErrorIsNil)
	s.PatchValue(series.DistroInfo, filename)

	expectedSeries := []string{"precise", "quantal", "raring", "saucy"}
	series := series.SupportedSeries()
	sort.Strings(series)
	c.Assert(series, gc.DeepEquals, expectedSeries)
}

func (s *supportedSeriesSuite) TestUpdateSeriesVersions(c *gc.C) {
	d := c.MkDir()
	filename := filepath.Join(d, "ubuntu.csv")
	err := ioutil.WriteFile(filename, []byte(distInfoData), 0644)
	c.Assert(err, jc.ErrorIsNil)
	s.PatchValue(series.DistroInfo, filename)

	expectedSeries := []string{"precise", "quantal", "raring", "saucy"}
	checkSeries := func() {
		series := series.SupportedSeries()
		sort.Strings(series)
		c.Assert(series, gc.DeepEquals, expectedSeries)
	}
	checkSeries()

	// Updating the file does not normally trigger an update;
	// we only refresh automatically one time. After that, we
	// must explicitly refresh.
	err = ioutil.WriteFile(filename, []byte(distInfoData2), 0644)
	c.Assert(err, jc.ErrorIsNil)
	checkSeries()

	expectedSeries = append([]string{"ornery"}, expectedSeries...)
	expectedSeries = append(expectedSeries, "trusty")
	series.UpdateSeriesVersions()
	checkSeries()
}

func (s *supportedSeriesSuite) TestOSSeries(c *gc.C) {
	cleanup := series.SetUbuntuSeries(make(map[string]string))
	defer cleanup()
	d := c.MkDir()
	filename := filepath.Join(d, "ubuntu.csv")
	err := ioutil.WriteFile(filename, []byte(distInfoData), 0644)
	c.Assert(err, jc.ErrorIsNil)
	s.PatchValue(series.DistroInfo, filename)

	osType, err := series.GetOSFromSeries("raring")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(osType, gc.Equals, os.Ubuntu)
}

const distInfoData = `version,codename,series,created,release,eol,eol-server
4.10,Warty Warthog,warty,2004-03-05,2004-10-20,2006-04-30
5.04,Hoary Hedgehog,hoary,2004-10-20,2005-04-08,2006-10-31
5.10,Breezy Badger,breezy,2005-04-08,2005-10-12,2007-04-13
6.06 LTS,Dapper Drake,dapper,2005-10-12,2006-06-01,2009-07-14,2011-06-01
6.10,Edgy Eft,edgy,2006-06-01,2006-10-26,2008-04-25
7.04,Feisty Fawn,feisty,2006-10-26,2007-04-19,2008-10-19
7.10,Gutsy Gibbon,gutsy,2007-04-19,2007-10-18,2009-04-18
8.04 LTS,Hardy Heron,hardy,2007-10-18,2008-04-24,2011-05-12,2013-05-09
8.10,Intrepid Ibex,intrepid,2008-04-24,2008-10-30,2010-04-30
9.04,Jaunty Jackalope,jaunty,2008-10-30,2009-04-23,2010-10-23
9.10,Karmic Koala,karmic,2009-04-23,2009-10-29,2011-04-29
10.04 LTS,Lucid Lynx,lucid,2009-10-29,2010-04-29,2013-05-09,2015-04-29
10.10,Maverick Meerkat,maverick,2010-04-29,2010-10-10,2012-04-10
11.04,Natty Narwhal,natty,2010-10-10,2011-04-28,2012-10-28
11.10,Oneiric Ocelot,oneiric,2011-04-28,2011-10-13,2013-05-09
12.04 LTS,Precise Pangolin,precise,2011-10-13,2012-04-26,2017-04-26
12.10,Quantal Quetzal,quantal,2012-04-26,2012-10-18,2014-04-18
13.04,Raring Ringtail,raring,2012-10-18,2013-04-25,2014-01-27
13.10,Saucy Salamander,saucy,2013-04-25,2013-10-17,2014-07-17
`

const distInfoData2 = distInfoData + `
14.04 LTS,Trusty Tahr,trusty,2013-10-17,2014-04-17,2019-04-17
94.04 LTS,Ornery Omega,ornery,2094-10-17,2094-04-17,2099-04-17
`
