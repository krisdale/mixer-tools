package swupd

import (
	"strings"
	"testing"
)

// Minversion support added in format 26
func TestFormats25to26Minversion(t *testing.T) {
	ts := newTestSwupd(t, "format25to26minversion")
	defer ts.cleanup()

	ts.Bundles = []string{"test-bundle"}

	// format25 MoM should NOT have minversion in header, which is introduced
	// in format26. (It should also not have it because minversion is set to 0)
	ts.Format = 25
	ts.addFile(10, "test-bundle", "/foo", "content")
	ts.createManifests(10)

	expSubs := []string{
		"MANIFEST\t25",
		"version:\t10",
		"previous:\t0",
		"filecount:\t2",
		"timestamp:\t",
		"contentsize:\t",
		"includes:\tos-core",
		"10\t/foo",
		"10\t/usr/share",
	}
	checkManifestContains(t, ts.Dir, "10", "test-bundle", expSubs...)

	nExpSubs := []string{
		"\t0\t/foo",
		".d..\t",
		"minversion:\t",
	}
	checkManifestNotContains(t, ts.Dir, "10", "test-bundle", nExpSubs...)

	// minversion now set to 20, but the MoM should still NOT have minversion
	// in header due to format25 being used
	ts.MinVersion = 20
	ts.addFile(20, "test-bundle", "/foo", "new content")
	ts.createManifests(20)

	expSubs = []string{
		"MANIFEST\t25",
		"version:\t20",
		"previous:\t10",
		"filecount:\t2",
		"includes:\tos-core",
		"20\t/foo",
	}
	checkManifestContains(t, ts.Dir, "20", "test-bundle", expSubs...)
	checkManifestNotContains(t, ts.Dir, "20", "MoM", "minversion:\t")

	// updated to format26, minversion still set to 20, so we should see
	// minversion  header in the MoM
	ts.Format = 26
	ts.addFile(30, "test-bundle", "/foo", "even newer content")
	ts.createManifests(30)
	expSubs = []string{
		"MANIFEST\t26",
		"version:\t30",
		"previous:\t20",
		"filecount:\t2",
		"includes:\tos-core",
	}
	checkManifestContains(t, ts.Dir, "30", "test-bundle", expSubs...)
	checkManifestContains(t, ts.Dir, "30", "MoM", "minversion:\t20")
}

// Iterative manifest support added in format 26
func TestFormats25to26IterativeManifest(t *testing.T) {
	ts := newTestSwupd(t, "format25to26iterativeManifest")
	defer ts.cleanup()

	ts.Bundles = []string{"test-bundle"}

	// Format 25 should not have iterative manifest support
	ts.Format = 25
	ts.addFile(10, "test-bundle", "/foo", "content")
	ts.createManifests(10)

	ts.addFile(20, "test-bundle", "/foo", "new content")
	ts.createManifests(20)
	checkManifestContains(t, ts.Dir, "20", "MoM", "MANIFEST\t25")

	// Iterative manifests should not have entries in the MoM or be generated
	checkManifestNotContains(t, ts.Dir, "20", "MoM", "I...\t")
	ts.checkNotExists("www/20/Manifest.test-bundle.I.10")
	ts.checkNotExists("www/20/Manifest.os-core.I.10")

	// Update to format26
	ts.Format = 26
	ts.addFile(30, "test-bundle", "/foo", "even newer content")
	ts.createManifests(30)

	ts.addFile(40, "test-bundle", "/foo", "more new content")
	ts.createManifests(40)
	checkManifestContains(t, ts.Dir, "40", "MoM", "MANIFEST\t26")

	// Updates in format 26 should support iterative manifests
	checkManifestContains(t, ts.Dir, "40", "MoM", "\ttest-bundle.I.30", "\tos-core.I.30")
	ts.checkExists("www/40/Manifest.test-bundle.I.30")
	ts.checkExists("www/40/Manifest.os-core.I.30")
}

// Delta manifest support added in format 26
func TestFormats25to26DeltaManifest(t *testing.T) {
	ts := newTestSwupd(t, "format25to26deltaManifest")
	defer ts.cleanup()

	ts.Bundles = []string{"test-bundle"}

	contents := strings.Repeat("large", 1000)
	if len(contents) < minimumSizeToMakeDeltaInBytes {
		t.Fatal("test content size is invalid")
	}

	// Format 25 should not have delta manifest support
	ts.Format = 25
	ts.addFile(10, "test-bundle", "/foo", contents+"A")
	ts.createManifests(10)

	ts.addFile(20, "test-bundle", "/foo", contents+"B")
	ts.createManifests(20)
	checkManifestContains(t, ts.Dir, "20", "MoM", "MANIFEST\t25")

	// Delta manifests should not exist
	ts.mustHashFile("image/10/full/foo")
	ts.mustHashFile("image/20/full/foo")
	ts.createPack("test-bundle", 10, 20, ts.path("image"))
	ts.checkNotExists("www/20/Manifest.test-bundle.D.10")

	// Update to format26
	ts.Format = 26
	ts.addFile(30, "test-bundle", "/foo", contents+"C")
	ts.createManifests(30)

	ts.addFile(40, "test-bundle", "/foo", contents+"D")
	ts.createManifests(40)
	checkManifestContains(t, ts.Dir, "40", "MoM", "MANIFEST\t26")

	// Delta manifests should be created in format 26
	ts.mustHashFile("image/30/full/foo")
	ts.mustHashFile("image/40/full/foo")
	ts.createPack("test-bundle", 30, 40, ts.path("image"))
	checkDeltaManifest(ts, 30, 40, "test-bundle", 1)
}

func TestFormat25BadContentSize(t *testing.T) {
	testCases := []struct {
		testName    string
		format      uint
		contentsize uint64
		expected    uint64
	}{
		// broken format
		{"format25: badMax + 1", 25, badMax + 1, badMax - 1},
		{"format25: badMax * 2", 25, badMax * 2, badMax - 1},
		{"format25: badMax", 25, badMax, badMax - 1},
		{"format25: badMax / 2", 25, badMax / 2, badMax / 2},
		// good format
		{"format26: badMax + 1", 26, badMax + 1, badMax + 1},
		{"format26: badMax * 2", 26, badMax * 2, badMax * 2},
		{"format26: badMax", 26, badMax, badMax},
		{"format26: badMax / 2", 26, badMax / 2, badMax / 2},
		// older good format
		{"format24: badMax + 1", 24, badMax + 1, badMax + 1},
		{"format24: badMax * 2", 24, badMax * 2, badMax * 2},
		{"format24: badMax", 24, badMax, badMax},
		{"format24: badMax / 2", 24, badMax / 2, badMax / 2},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			m := &Manifest{
				Header: ManifestHeader{
					Format:      tc.format,
					ContentSize: tc.contentsize,
				},
			}
			m.setMaxContentSizeForFormat()
			if m.Header.ContentSize != tc.expected {
				t.Errorf("%d contentsize set to %d, expected %d",
					tc.contentsize,
					m.Header.ContentSize,
					tc.expected,
				)
			}
		})
	}
}

// Experimental bundles added in format 27
func TestFormats26to27ExperimentalBundles(t *testing.T) {
	ts := newTestSwupd(t, "format26to27ExperimentalBundles")
	defer ts.cleanup()

	var header BundleHeader
	header.Status = "Experimental"

	// Format 26 should not recognize experimental bundles
	ts.Format = 26
	ts.Bundles = []string{"test-bundle1"}
	ts.addFile(10, "test-bundle1", "/foo", "content")
	ts.addHeader(10, "test-bundle1", header)
	ts.createManifests(10)
	checkManifestNotContains(t, ts.Dir, "10", "MoM", "Me..\t")

	// Format 27 should recognize experimental bundles
	ts.Format = 27
	ts.Bundles = []string{"test-bundle2"}
	ts.addFile(20, "test-bundle2", "/foo", "content")
	ts.addHeader(20, "test-bundle2", header)
	ts.createManifests(20)
	checkManifestContains(t, ts.Dir, "20", "MoM", "Me..\t")
}
