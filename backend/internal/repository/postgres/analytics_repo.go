// File: analytics_repo.go
// Purpose: Implement PostgreSQL repository methods for visitor analytics ingest and queries.
// Module: backend/internal/repository/postgres, analytics persistence layer.
// Related: repository.go, service analytics flows, and admin/public analytics handlers.
package postgres

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
)

const analyticsVisitSelectColumns = `
SELECT ve.id::text, ve.event_id, ve.event_type, ve.occurred_at, ve.received_at, ve.visitor_id, ve.session_id, ve.path,
       COALESCE(ve.route_name,''), COALESCE(ve.page_title,''), COALESCE(ve.referrer,''), COALESCE(ve.referrer_host,''),
       COALESCE(ve.content_type,''), COALESCE(ve.content_id,''), COALESCE(ve.content_slug,''), COALESCE(ve.locale,''),
       COALESCE(ve.screen_width,0), COALESCE(ve.screen_height,0), COALESCE(ve.viewport_width,0), COALESCE(ve.viewport_height,0),
       COALESCE(ve.timezone,''), ve.ip, COALESCE(ipf.country_code,''), COALESCE(ipf.country_name,''), COALESCE(ipf.region_name,''),
       COALESCE(ipf.city_name,''), COALESCE(ipf.isp,''), COALESCE(ve.user_agent,''), COALESCE(ve.device_type,''), COALESCE(ve.browser_name,''),
       COALESCE(ve.os_name,''), ve.is_bot, ve.created_at
FROM visit_events ve
LEFT JOIN ip_profiles ipf ON ipf.ip = ve.ip`

func (r *Repository) GetIPProfile(ip string) (domain.IPProfile, bool, error) {
	var item domain.IPProfile
	err := r.db.QueryRow(`
SELECT ip, COALESCE(country_code,''), COALESCE(country_name,''), COALESCE(region_name,''), COALESCE(city_name,''),
       COALESCE(isp,''), COALESCE(raw_region,''), COALESCE(source,''), resolved_at, created_at, updated_at
FROM ip_profiles
WHERE ip = $1
`, strings.TrimSpace(ip)).Scan(
		&item.IP,
		&item.CountryCode,
		&item.CountryName,
		&item.RegionName,
		&item.CityName,
		&item.ISP,
		&item.RawRegion,
		&item.Source,
		&item.ResolvedAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.IPProfile{}, false, nil
		}
		return domain.IPProfile{}, false, err
	}
	return item, true, nil
}

func (r *Repository) UpsertIPProfile(profile domain.IPProfile) (domain.IPProfile, error) {
	err := r.db.QueryRow(`
INSERT INTO ip_profiles (
    ip, country_code, country_name, region_name, city_name, isp, raw_region, source, resolved_at
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
ON CONFLICT (ip) DO UPDATE
SET country_code = EXCLUDED.country_code,
    country_name = EXCLUDED.country_name,
    region_name = EXCLUDED.region_name,
    city_name = EXCLUDED.city_name,
    isp = EXCLUDED.isp,
    raw_region = EXCLUDED.raw_region,
    source = EXCLUDED.source,
    resolved_at = EXCLUDED.resolved_at,
    updated_at = NOW()
RETURNING ip, COALESCE(country_code,''), COALESCE(country_name,''), COALESCE(region_name,''), COALESCE(city_name,''),
          COALESCE(isp,''), COALESCE(raw_region,''), COALESCE(source,''), resolved_at, created_at, updated_at
`,
		strings.TrimSpace(profile.IP),
		nullableString(profile.CountryCode),
		nullableString(profile.CountryName),
		nullableString(profile.RegionName),
		nullableString(profile.CityName),
		nullableString(profile.ISP),
		nullableString(profile.RawRegion),
		nullableString(profile.Source),
		profile.ResolvedAt.UTC(),
	).Scan(
		&profile.IP,
		&profile.CountryCode,
		&profile.CountryName,
		&profile.RegionName,
		&profile.CityName,
		&profile.ISP,
		&profile.RawRegion,
		&profile.Source,
		&profile.ResolvedAt,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	return profile, err
}

func (r *Repository) CreateVisitEvents(events []domain.VisitEvent) (domain.AnalyticsIngestResult, error) {
	result := domain.AnalyticsIngestResult{}
	if len(events) == 0 {
		return result, nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return result, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare(`
INSERT INTO visit_events (
    event_id, event_type, occurred_at, visitor_id, session_id, path, route_name, page_title, referrer, referrer_host,
    content_type, content_id, content_slug, locale, screen_width, screen_height, viewport_width, viewport_height,
    timezone, ip, user_agent, device_type, browser_name, os_name, is_bot
)
VALUES (
    $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,
    $11,$12,$13,$14,$15,$16,$17,$18,
    $19,$20,$21,$22,$23,$24,$25
)
ON CONFLICT (event_id) DO NOTHING
`)
	if err != nil {
		return result, err
	}
	defer stmt.Close()

	for _, event := range events {
		execResult, execErr := stmt.Exec(
			event.EventID,
			event.EventType,
			event.OccurredAt.UTC(),
			event.VisitorID,
			event.SessionID,
			event.Path,
			nullableString(event.RouteName),
			nullableString(event.PageTitle),
			nullableString(event.Referrer),
			nullableString(event.ReferrerHost),
			nullableString(event.ContentType),
			nullableString(event.ContentID),
			nullableString(event.ContentSlug),
			nullableString(event.Locale),
			nullableInt(event.ScreenWidth),
			nullableInt(event.ScreenHeight),
			nullableInt(event.ViewportWidth),
			nullableInt(event.ViewportHeight),
			nullableString(event.Timezone),
			event.IP,
			nullableString(event.UserAgent),
			nullableString(event.DeviceType),
			nullableString(event.BrowserName),
			nullableString(event.OSName),
			event.IsBot,
		)
		if execErr != nil {
			err = execErr
			return result, err
		}
		affected, _ := execResult.RowsAffected()
		if affected > 0 {
			result.Accepted++
		} else {
			result.Deduplicated++
		}
	}

	err = tx.Commit()
	return result, err
}

func (r *Repository) GetAnalyticsOverview(days int) (domain.AnalyticsOverview, error) {
	start, end := analyticsRange(days)
	overview := domain.AnalyticsOverview{
		RangeStart:      start,
		RangeEnd:        end.Add(-time.Nanosecond),
		TopPaths:        []domain.AnalyticsPathStat{},
		TopReferrers:    []domain.AnalyticsReferrerStat{},
		DeviceBreakdown: []domain.AnalyticsDeviceStat{},
		Daily:           []domain.AnalyticsDailyStat{},
	}

	err := r.db.QueryRow(`
SELECT COUNT(*), COUNT(DISTINCT visitor_id), COUNT(DISTINCT ip), COUNT(DISTINCT session_id)
FROM visit_events
WHERE event_type='page_view' AND occurred_at >= $1 AND occurred_at < $2
`, start, end).Scan(&overview.PageViews, &overview.UniqueVisitors, &overview.UniqueIPs, &overview.UniqueSessions)
	if err != nil {
		return overview, err
	}

	topPathRows, err := r.db.Query(`
SELECT path, COALESCE(content_type,''), COALESCE(content_id,''), COALESCE(content_slug,''),
       COUNT(*) AS page_views, COUNT(DISTINCT visitor_id) AS unique_visitors, COUNT(DISTINCT ip) AS unique_ips,
       MAX(occurred_at) AS last_visited_at
FROM visit_events
WHERE event_type='page_view' AND occurred_at >= $1 AND occurred_at < $2
GROUP BY path, content_type, content_id, content_slug
ORDER BY page_views DESC, last_visited_at DESC
LIMIT 10
`, start, end)
	if err != nil {
		return overview, err
	}
	defer topPathRows.Close()
	for topPathRows.Next() {
		var item domain.AnalyticsPathStat
		if scanErr := topPathRows.Scan(&item.Path, &item.ContentType, &item.ContentID, &item.ContentSlug, &item.PageViews, &item.UniqueVisitors, &item.UniqueIPs, &item.LastVisitedAt); scanErr == nil {
			overview.TopPaths = append(overview.TopPaths, item)
		}
	}

	refRows, err := r.db.Query(`
SELECT referrer_host, COUNT(*) AS visits
FROM visit_events
WHERE event_type='page_view' AND occurred_at >= $1 AND occurred_at < $2 AND referrer_host IS NOT NULL AND referrer_host <> ''
GROUP BY referrer_host
ORDER BY visits DESC, referrer_host ASC
LIMIT 10
`, start, end)
	if err != nil {
		return overview, err
	}
	defer refRows.Close()
	for refRows.Next() {
		var item domain.AnalyticsReferrerStat
		if scanErr := refRows.Scan(&item.ReferrerHost, &item.Visits); scanErr == nil {
			overview.TopReferrers = append(overview.TopReferrers, item)
		}
	}

	deviceRows, err := r.db.Query(`
SELECT COALESCE(device_type,'unknown') AS device_type, COUNT(*) AS visits
FROM visit_events
WHERE event_type='page_view' AND occurred_at >= $1 AND occurred_at < $2
GROUP BY COALESCE(device_type,'unknown')
ORDER BY visits DESC, device_type ASC
LIMIT 10
`, start, end)
	if err != nil {
		return overview, err
	}
	defer deviceRows.Close()
	for deviceRows.Next() {
		var item domain.AnalyticsDeviceStat
		if scanErr := deviceRows.Scan(&item.DeviceType, &item.Visits); scanErr == nil {
			overview.DeviceBreakdown = append(overview.DeviceBreakdown, item)
		}
	}

	dailyRows, err := r.db.Query(`
WITH days AS (
  SELECT generate_series($1::date, ($2::date - 1), interval '1 day')::date AS day
)
SELECT to_char(days.day, 'YYYY-MM-DD') AS date,
       COALESCE(COUNT(ve.id), 0) AS page_views,
       COALESCE(COUNT(DISTINCT ve.visitor_id), 0) AS unique_visitors,
       COALESCE(COUNT(DISTINCT ve.ip), 0) AS unique_ips
FROM days
LEFT JOIN visit_events ve
  ON ve.event_type='page_view'
 AND ve.occurred_at >= days.day
 AND ve.occurred_at < days.day + interval '1 day'
GROUP BY days.day
ORDER BY days.day ASC
`, start, end)
	if err != nil {
		return overview, err
	}
	defer dailyRows.Close()
	for dailyRows.Next() {
		var item domain.AnalyticsDailyStat
		if scanErr := dailyRows.Scan(&item.Date, &item.PageViews, &item.UniqueVisitors, &item.UniqueIPs); scanErr == nil {
			overview.Daily = append(overview.Daily, item)
		}
	}

	return overview, nil
}

func (r *Repository) ListAnalyticsPages(page, pageSize, days int, path, contentType string) ([]domain.AnalyticsPathStat, int, error) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize
	start, end := analyticsRange(days)

	whereClause, args := analyticsWhereClause(start, end, path, "", "", "", contentType, "", "", "", "", "", true)

	countQuery := `
SELECT COUNT(*) FROM (
  SELECT 1
  FROM visit_events
  WHERE ` + whereClause + `
  GROUP BY path, content_type, content_id, content_slug
) grouped`
	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return []domain.AnalyticsPathStat{}, 0, err
	}

	queryArgs := append(args, pageSize, offset)
	query := `
SELECT path, COALESCE(content_type,''), COALESCE(content_id,''), COALESCE(content_slug,''),
       COUNT(*) AS page_views, COUNT(DISTINCT visitor_id) AS unique_visitors, COUNT(DISTINCT ip) AS unique_ips,
       MAX(occurred_at) AS last_visited_at
FROM visit_events
WHERE ` + whereClause + `
GROUP BY path, content_type, content_id, content_slug
ORDER BY page_views DESC, last_visited_at DESC
LIMIT $` + strconv.Itoa(len(queryArgs)-1) + ` OFFSET $` + strconv.Itoa(len(queryArgs))
	rows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		return []domain.AnalyticsPathStat{}, 0, err
	}
	defer rows.Close()

	items := make([]domain.AnalyticsPathStat, 0)
	for rows.Next() {
		var item domain.AnalyticsPathStat
		if scanErr := rows.Scan(&item.Path, &item.ContentType, &item.ContentID, &item.ContentSlug, &item.PageViews, &item.UniqueVisitors, &item.UniqueIPs, &item.LastVisitedAt); scanErr == nil {
			items = append(items, item)
		}
	}
	return items, total, nil
}

func (r *Repository) ListAnalyticsVisits(page, pageSize, days int, path, eventType, visitorID, sessionID, contentType, ip, deviceType, browserName, osName, isBot, countryName, regionName, cityName, isp string) ([]domain.VisitEvent, int, error) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize
	start, end := analyticsRange(days)

	whereClause, args := analyticsVisitWhereClause(start, end, path, eventType, visitorID, sessionID, contentType, ip, deviceType, browserName, osName, isBot, countryName, regionName, cityName, isp)

	var total int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM visit_events ve LEFT JOIN ip_profiles ipf ON ipf.ip = ve.ip WHERE `+whereClause, args...).Scan(&total); err != nil {
		return []domain.VisitEvent{}, 0, err
	}

	queryArgs := append(args, pageSize, offset)
	query := analyticsVisitSelectColumns + `
WHERE ` + whereClause + `
ORDER BY ve.occurred_at DESC
LIMIT $` + strconv.Itoa(len(queryArgs)-1) + ` OFFSET $` + strconv.Itoa(len(queryArgs))
	rows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		return []domain.VisitEvent{}, 0, err
	}
	return scanVisitEventRows(rows), total, nil
}

func analyticsRange(days int) (time.Time, time.Time) {
	now := time.Now().UTC()
	end := now.Truncate(24 * time.Hour).Add(24 * time.Hour)
	start := end.AddDate(0, 0, -days)
	return start, end
}

func analyticsVisitWhereClause(start, end time.Time, path, eventType, visitorID, sessionID, contentType, ip, deviceType, browserName, osName, isBot, countryName, regionName, cityName, isp string) (string, []any) {
	conditions := []string{"ve.occurred_at >= $1", "ve.occurred_at < $2"}
	args := []any{start, end}
	if strings.TrimSpace(path) != "" {
		args = append(args, "%"+strings.TrimSpace(path)+"%")
		conditions = append(conditions, "ve.path ILIKE $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(eventType) != "" {
		args = append(args, strings.TrimSpace(eventType))
		conditions = append(conditions, "ve.event_type = $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(visitorID) != "" {
		args = append(args, strings.TrimSpace(visitorID))
		conditions = append(conditions, "ve.visitor_id = $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(sessionID) != "" {
		args = append(args, strings.TrimSpace(sessionID))
		conditions = append(conditions, "ve.session_id = $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(contentType) != "" {
		args = append(args, strings.TrimSpace(contentType))
		conditions = append(conditions, "ve.content_type = $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(ip) != "" {
		args = append(args, strings.TrimSpace(ip))
		conditions = append(conditions, "ve.ip = $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(deviceType) != "" {
		args = append(args, strings.TrimSpace(deviceType))
		conditions = append(conditions, "ve.device_type = $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(browserName) != "" {
		args = append(args, strings.TrimSpace(browserName))
		conditions = append(conditions, "ve.browser_name = $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(osName) != "" {
		args = append(args, strings.TrimSpace(osName))
		conditions = append(conditions, "ve.os_name = $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(countryName) != "" {
		args = append(args, strings.TrimSpace(countryName))
		conditions = append(conditions, "ipf.country_name = $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(regionName) != "" {
		args = append(args, strings.TrimSpace(regionName))
		conditions = append(conditions, "ipf.region_name = $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(cityName) != "" {
		args = append(args, strings.TrimSpace(cityName))
		conditions = append(conditions, "ipf.city_name = $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(isp) != "" {
		args = append(args, strings.TrimSpace(isp))
		conditions = append(conditions, "ipf.isp = $"+strconv.Itoa(len(args)))
	}
	switch strings.ToLower(strings.TrimSpace(isBot)) {
	case "true", "1", "yes":
		conditions = append(conditions, "ve.is_bot = TRUE")
	case "false", "0", "no":
		conditions = append(conditions, "ve.is_bot = FALSE")
	}
	return strings.Join(conditions, " AND "), args
}

func analyticsWhereClause(start, end time.Time, path, eventType, visitorID, sessionID, contentType, ip, deviceType, browserName, osName, isBot string, pageViewsOnly bool) (string, []any) {
	conditions := []string{"occurred_at >= $1", "occurred_at < $2"}
	args := []any{start, end}
	if pageViewsOnly {
		conditions = append(conditions, "event_type='page_view'")
	}
	if strings.TrimSpace(path) != "" {
		args = append(args, "%"+strings.TrimSpace(path)+"%")
		conditions = append(conditions, "path ILIKE $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(eventType) != "" {
		args = append(args, strings.TrimSpace(eventType))
		conditions = append(conditions, "event_type = $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(visitorID) != "" {
		args = append(args, strings.TrimSpace(visitorID))
		conditions = append(conditions, "visitor_id = $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(sessionID) != "" {
		args = append(args, strings.TrimSpace(sessionID))
		conditions = append(conditions, "session_id = $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(contentType) != "" {
		args = append(args, strings.TrimSpace(contentType))
		conditions = append(conditions, "content_type = $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(ip) != "" {
		args = append(args, strings.TrimSpace(ip))
		conditions = append(conditions, "ip = $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(deviceType) != "" {
		args = append(args, strings.TrimSpace(deviceType))
		conditions = append(conditions, "device_type = $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(browserName) != "" {
		args = append(args, strings.TrimSpace(browserName))
		conditions = append(conditions, "browser_name = $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(osName) != "" {
		args = append(args, strings.TrimSpace(osName))
		conditions = append(conditions, "os_name = $"+strconv.Itoa(len(args)))
	}
	switch strings.ToLower(strings.TrimSpace(isBot)) {
	case "true", "1", "yes":
		conditions = append(conditions, "is_bot = TRUE")
	case "false", "0", "no":
		conditions = append(conditions, "is_bot = FALSE")
	}
	return strings.Join(conditions, " AND "), args
}

func nullableInt(v int) any {
	if v <= 0 {
		return nil
	}
	return v
}

func scanVisitEventRows(rows *sql.Rows) []domain.VisitEvent {
	defer rows.Close()

	items := make([]domain.VisitEvent, 0)
	for rows.Next() {
		var item domain.VisitEvent
		if err := rows.Scan(
			&item.ID,
			&item.EventID,
			&item.EventType,
			&item.OccurredAt,
			&item.ReceivedAt,
			&item.VisitorID,
			&item.SessionID,
			&item.Path,
			&item.RouteName,
			&item.PageTitle,
			&item.Referrer,
			&item.ReferrerHost,
			&item.ContentType,
			&item.ContentID,
			&item.ContentSlug,
			&item.Locale,
			&item.ScreenWidth,
			&item.ScreenHeight,
			&item.ViewportWidth,
			&item.ViewportHeight,
			&item.Timezone,
			&item.IP,
			&item.CountryCode,
			&item.CountryName,
			&item.RegionName,
			&item.CityName,
			&item.ISP,
			&item.UserAgent,
			&item.DeviceType,
			&item.BrowserName,
			&item.OSName,
			&item.IsBot,
			&item.CreatedAt,
		); err == nil {
			items = append(items, item)
		}
	}
	return items
}
