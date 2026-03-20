// File: ip2region.go
// Purpose: Resolve visitor IP addresses into offline geographic metadata through ip2region xdb files.
// Module: backend/internal/ipgeo, infrastructure enrichment layer.
// Related: config analytics settings, content analytics service, and ip_profiles persistence.
package ipgeo

import (
	"net"
	"strings"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

type IP2RegionResolver struct {
	ipv4DBPath string
	ipv6DBPath string
}

func NewIP2RegionResolver(ipv4DBPath, ipv6DBPath string) *IP2RegionResolver {
	return &IP2RegionResolver{
		ipv4DBPath: strings.TrimSpace(ipv4DBPath),
		ipv6DBPath: strings.TrimSpace(ipv6DBPath),
	}
}

func (r *IP2RegionResolver) Lookup(ip string) (domain.IPProfile, bool, error) {
	normalized := normalizeLookupIP(ip)
	if normalized == "" {
		return domain.IPProfile{}, false, nil
	}

	version, err := xdb.VersionFromIP(normalized)
	if err != nil {
		return domain.IPProfile{}, false, err
	}

	dbPath := r.dbPathForVersion(version)
	if dbPath == "" {
		return domain.IPProfile{}, false, nil
	}

	searcher, err := xdb.NewWithFileOnly(version, dbPath)
	if err != nil {
		return domain.IPProfile{}, false, err
	}
	defer searcher.Close()

	rawRegion, err := searcher.SearchByStr(normalized)
	if err != nil {
		return domain.IPProfile{}, false, err
	}

	profile := domain.IPProfile{
		IP:         normalized,
		RawRegion:  strings.TrimSpace(rawRegion),
		Source:     "ip2region",
		ResolvedAt: time.Now().UTC(),
	}
	parts := strings.Split(profile.RawRegion, "|")
	if len(parts) > 0 {
		profile.CountryName = normalizeRegionPart(parts[0])
	}
	if len(parts) > 1 {
		profile.RegionName = normalizeRegionPart(parts[1])
	}
	if len(parts) > 2 {
		profile.CityName = normalizeRegionPart(parts[2])
	}
	if len(parts) > 3 {
		profile.ISP = normalizeRegionPart(parts[3])
	}
	if len(parts) > 4 {
		profile.CountryCode = strings.ToUpper(normalizeRegionPart(parts[4]))
	}

	if profile.RawRegion == "" && profile.CountryName == "" && profile.RegionName == "" && profile.CityName == "" && profile.ISP == "" {
		return domain.IPProfile{}, false, nil
	}
	return profile, true, nil
}

func (r *IP2RegionResolver) dbPathForVersion(version *xdb.Version) string {
	if version == nil {
		return ""
	}
	switch version.Id {
	case xdb.IPv4VersionNo:
		return r.ipv4DBPath
	case xdb.IPv6VersionNo:
		return r.ipv6DBPath
	default:
		return ""
	}
}

func normalizeLookupIP(raw string) string {
	parsed := net.ParseIP(strings.TrimSpace(raw))
	if parsed == nil {
		return ""
	}
	if parsed.IsLoopback() || parsed.IsPrivate() || parsed.IsUnspecified() || parsed.IsMulticast() || parsed.IsLinkLocalUnicast() || parsed.IsLinkLocalMulticast() {
		return ""
	}
	return parsed.String()
}

func normalizeRegionPart(raw string) string {
	value := strings.TrimSpace(raw)
	if value == "" || value == "0" {
		return ""
	}
	return value
}
