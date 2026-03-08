// File: site_repo.go
// Purpose: Implement PostgreSQL repository methods for site settings, nav, social, footer, and slots.
// Module: backend/internal/repository/postgres, site persistence layer.
// Related: helpers.go and service site flows.
package postgres

import (
	"encoding/json"

	"github.com/anxcye/ancy-blog/backend/internal/apperr"
	"github.com/anxcye/ancy-blog/backend/internal/domain"
)

func (r *Repository) GetSiteSettings() domain.SiteSettings {
	var s domain.SiteSettings
	err := r.db.QueryRow(`
SELECT site_name, COALESCE(avatar_url,''), COALESCE(favicon_url,''), COALESCE(hero_intro_md,''), default_locale,
       comment_enabled, comment_require_approval, link_submission_enabled,
       COALESCE(site_description,''), COALESCE(seo_keywords,''), COALESCE(og_image_url,'')
FROM site_settings ORDER BY created_at ASC LIMIT 1`).
		Scan(&s.SiteName, &s.AvatarURL, &s.FaviconURL, &s.HeroIntroMD, &s.DefaultLocale,
			&s.CommentEnabled, &s.CommentRequireApproval, &s.LinkSubmissionEnabled,
			&s.SiteDescription, &s.SeoKeywords, &s.OgImageURL)
	if err != nil {
		return domain.SiteSettings{SiteName: "Ancy Blog", DefaultLocale: "en", CommentEnabled: true, LinkSubmissionEnabled: true}
	}
	return s
}

func (r *Repository) UpdateSiteSettings(settings domain.SiteSettings) domain.SiteSettings {
	var id string
	err := r.db.QueryRow(`SELECT id::text FROM site_settings ORDER BY created_at ASC LIMIT 1`).Scan(&id)
	if err != nil {
		_ = r.db.QueryRow(`
INSERT INTO site_settings (site_name, avatar_url, hero_intro_md, default_locale,
    favicon_url, comment_enabled, comment_require_approval, link_submission_enabled, site_description, seo_keywords, og_image_url)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id::text
`, settings.SiteName, nullableString(settings.AvatarURL), nullableString(settings.FaviconURL), nullableString(settings.HeroIntroMD), settings.DefaultLocale,
			settings.CommentEnabled, settings.CommentRequireApproval, settings.LinkSubmissionEnabled,
			nullableString(settings.SiteDescription), nullableString(settings.SeoKeywords), nullableString(settings.OgImageURL)).Scan(&id)
		return settings
	}
	_, _ = r.db.Exec(`
UPDATE site_settings
SET site_name=$2, avatar_url=$3, favicon_url=$4, hero_intro_md=$5, default_locale=$6,
    comment_enabled=$7, comment_require_approval=$8, link_submission_enabled=$9,
    site_description=$10, seo_keywords=$11, og_image_url=$12, updated_at=NOW()
WHERE id=$1
`, id, settings.SiteName, nullableString(settings.AvatarURL), nullableString(settings.FaviconURL), nullableString(settings.HeroIntroMD), settings.DefaultLocale,
		settings.CommentEnabled, settings.CommentRequireApproval, settings.LinkSubmissionEnabled,
		nullableString(settings.SiteDescription), nullableString(settings.SeoKeywords), nullableString(settings.OgImageURL))
	return settings
}

func (r *Repository) GetAdminPasswordHash() (string, bool) {
	var hash string
	err := r.db.QueryRow(`SELECT admin_password_hash FROM site_settings WHERE admin_password_hash IS NOT NULL ORDER BY created_at ASC LIMIT 1`).Scan(&hash)
	if err != nil || hash == "" {
		return "", false
	}
	return hash, true
}

func (r *Repository) SetAdminPasswordHash(hash string) error {
	var id string
	err := r.db.QueryRow(`SELECT id::text FROM site_settings ORDER BY created_at ASC LIMIT 1`).Scan(&id)
	if err != nil {
		// No row yet — create default row with hash
		return r.db.QueryRow(`
INSERT INTO site_settings (site_name, default_locale, admin_password_hash)
VALUES ('Ancy Blog','en',$1) RETURNING id::text
`, hash).Scan(&id)
	}
	_, err = r.db.Exec(`UPDATE site_settings SET admin_password_hash=$2, updated_at=NOW() WHERE id=$1`, id, hash)
	return err
}

func (r *Repository) GetTranslationPolicy() domain.TranslationPolicy {
	var raw []byte
	err := r.db.QueryRow(`SELECT translation_policy FROM site_settings ORDER BY created_at ASC LIMIT 1`).Scan(&raw)
	if err != nil || raw == nil {
		return domain.TranslationPolicy{TargetLocales: []string{}}
	}
	var policy domain.TranslationPolicy
	if err := json.Unmarshal(raw, &policy); err != nil {
		return domain.TranslationPolicy{TargetLocales: []string{}}
	}
	if policy.TargetLocales == nil {
		policy.TargetLocales = []string{}
	}
	return policy
}

func (r *Repository) UpdateTranslationPolicy(policy domain.TranslationPolicy) error {
	raw, err := json.Marshal(policy)
	if err != nil {
		return err
	}
	var id string
	err = r.db.QueryRow(`SELECT id::text FROM site_settings ORDER BY created_at ASC LIMIT 1`).Scan(&id)
	if err != nil {
		// No row yet — create a default row with just the policy
		return r.db.QueryRow(`
INSERT INTO site_settings (site_name, default_locale, translation_policy)
VALUES ('Ancy Blog','en',$1) RETURNING id::text
`, raw).Scan(&id)
	}
	_, err = r.db.Exec(`UPDATE site_settings SET translation_policy=$2, updated_at=NOW() WHERE id=$1`, id, raw)
	return err
}

func (r *Repository) CreateFooterItem(item domain.FooterItem) (domain.FooterItem, error) {
	var id string
	err := r.db.QueryRow(`
INSERT INTO footer_items (label, link_type, internal_article_slug, external_url, row_num, order_num, enabled)
VALUES ($1,$2,$3,$4,$5,$6,$7)
RETURNING id::text
`, item.Label, item.LinkType, nullableString(item.InternalArticleSlug), nullableString(item.ExternalURL), item.RowNum, item.OrderNum, item.Enabled).Scan(&id)
	if err != nil {
		return domain.FooterItem{}, err
	}
	item.ID = id
	return item, nil
}

func (r *Repository) UpdateFooterItem(id string, item domain.FooterItem) (domain.FooterItem, error) {
	res, err := r.db.Exec(`
UPDATE footer_items
SET label=$2, link_type=$3, internal_article_slug=$4, external_url=$5, row_num=$6, order_num=$7, enabled=$8, updated_at=NOW()
WHERE id=$1 AND deleted_at IS NULL
`, id, item.Label, item.LinkType, nullableString(item.InternalArticleSlug), nullableString(item.ExternalURL), item.RowNum, item.OrderNum, item.Enabled)
	if err != nil {
		return domain.FooterItem{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.FooterItem{}, apperr.ErrFooterItemNotFound
	}
	item.ID = id
	return item, nil
}

func (r *Repository) DeleteFooterItem(id string) bool {
	res, err := r.db.Exec(`UPDATE footer_items SET deleted_at=NOW(), updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return false
	}
	n, _ := res.RowsAffected()
	return n > 0
}

func (r *Repository) ListFooterItems() []domain.FooterItem {
	rows, err := r.db.Query(`
SELECT id::text, label, link_type, COALESCE(internal_article_slug,''), COALESCE(external_url,''), row_num, order_num, enabled
FROM footer_items
WHERE enabled=TRUE AND deleted_at IS NULL
ORDER BY row_num ASC, order_num ASC
`)
	if err != nil {
		return []domain.FooterItem{}
	}
	defer rows.Close()
	items := make([]domain.FooterItem, 0)
	for rows.Next() {
		var f domain.FooterItem
		if err := rows.Scan(&f.ID, &f.Label, &f.LinkType, &f.InternalArticleSlug, &f.ExternalURL, &f.RowNum, &f.OrderNum, &f.Enabled); err == nil {
			items = append(items, f)
		}
	}
	return items
}

func (r *Repository) CreateSocialLink(item domain.SocialLink) (domain.SocialLink, error) {
	var id string
	err := r.db.QueryRow(`
INSERT INTO social_links (platform, title, url, icon_key, order_num, enabled)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING id::text
`, item.Platform, item.Title, item.URL, nullableString(item.IconKey), item.OrderNum, item.Enabled).Scan(&id)
	if err != nil {
		return domain.SocialLink{}, err
	}
	item.ID = id
	return item, nil
}

func (r *Repository) UpdateSocialLink(id string, item domain.SocialLink) (domain.SocialLink, error) {
	res, err := r.db.Exec(`
UPDATE social_links
SET platform=$2, title=$3, url=$4, icon_key=$5, order_num=$6, enabled=$7, updated_at=NOW()
WHERE id=$1 AND deleted_at IS NULL
`, id, item.Platform, item.Title, item.URL, nullableString(item.IconKey), item.OrderNum, item.Enabled)
	if err != nil {
		return domain.SocialLink{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.SocialLink{}, apperr.ErrSocialLinkNotFound
	}
	item.ID = id
	return item, nil
}

func (r *Repository) DeleteSocialLink(id string) bool {
	res, err := r.db.Exec(`UPDATE social_links SET deleted_at=NOW(), updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return false
	}
	n, _ := res.RowsAffected()
	return n > 0
}

func (r *Repository) ListSocialLinks() []domain.SocialLink {
	rows, err := r.db.Query(`
SELECT id::text, platform, title, url, COALESCE(icon_key,''), order_num, enabled
FROM social_links
WHERE enabled=TRUE AND deleted_at IS NULL
ORDER BY order_num ASC
`)
	if err != nil {
		return []domain.SocialLink{}
	}
	defer rows.Close()
	items := make([]domain.SocialLink, 0)
	for rows.Next() {
		var s domain.SocialLink
		if err := rows.Scan(&s.ID, &s.Platform, &s.Title, &s.URL, &s.IconKey, &s.OrderNum, &s.Enabled); err == nil {
			items = append(items, s)
		}
	}
	return items
}

func (r *Repository) CreateNavItem(item domain.NavItem) (domain.NavItem, error) {
	var id string
	err := r.db.QueryRow(`
INSERT INTO nav_items (name, key, type, target_type, target_value, order_num, enabled, parent_id)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
RETURNING id::text
`, item.Name, item.Key, item.Type, item.TargetType, nullableString(item.TargetValue), item.OrderNum, item.Enabled, nullableString(item.ParentID)).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.NavItem{}, apperr.ErrValidation
		}
		return domain.NavItem{}, err
	}
	item.ID = id
	return item, nil
}

func (r *Repository) UpdateNavItem(id string, item domain.NavItem) (domain.NavItem, error) {
	res, err := r.db.Exec(`
UPDATE nav_items
SET name=$2, key=$3, type=$4, target_type=$5, target_value=$6, order_num=$7, enabled=$8, parent_id=$9, updated_at=NOW()
WHERE id=$1 AND deleted_at IS NULL
`, id, item.Name, item.Key, item.Type, item.TargetType, nullableString(item.TargetValue), item.OrderNum, item.Enabled, nullableString(item.ParentID))
	if err != nil {
		return domain.NavItem{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.NavItem{}, apperr.ErrNavItemNotFound
	}
	item.ID = id
	return item, nil
}

func (r *Repository) DeleteNavItem(id string) bool {
	res, err := r.db.Exec(`UPDATE nav_items SET deleted_at=NOW(), updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return false
	}
	n, _ := res.RowsAffected()
	return n > 0
}

func (r *Repository) ListNavItems() []domain.NavItem {
	rows, err := r.db.Query(`
SELECT id::text, COALESCE(parent_id::text,''), name, key, type, target_type, COALESCE(target_value,''), order_num, enabled
FROM nav_items
WHERE enabled=TRUE AND deleted_at IS NULL
ORDER BY order_num ASC
`)
	if err != nil {
		return []domain.NavItem{}
	}
	defer rows.Close()
	// flat scan first
	flat := make([]domain.NavItem, 0)
	for rows.Next() {
		var n domain.NavItem
		if err := rows.Scan(&n.ID, &n.ParentID, &n.Name, &n.Key, &n.Type, &n.TargetType, &n.TargetValue, &n.OrderNum, &n.Enabled); err == nil {
			flat = append(flat, n)
		}
	}
	// build tree
	byID := make(map[string]*domain.NavItem, len(flat))
	for i := range flat {
		byID[flat[i].ID] = &flat[i]
	}
	result := make([]domain.NavItem, 0)
	for i := range flat {
		item := &flat[i]
		if item.ParentID == "" {
			result = append(result, *item)
		} else if parent, ok := byID[item.ParentID]; ok {
			parent.Children = append(parent.Children, *item)
		}
	}
	// flush updated parents into result
	for i, r2 := range result {
		if p, ok := byID[r2.ID]; ok {
			result[i] = *p
		}
	}
	return result
}

func (r *Repository) CreateContentSlot(slot domain.ContentSlot) (domain.ContentSlot, error) {
	var id string
	err := r.db.QueryRow(`
INSERT INTO content_slots (slot_key, name, description, enabled)
VALUES ($1,$2,$3,$4)
RETURNING id::text
`, slot.SlotKey, slot.Name, nullableString(slot.Description), slot.Enabled).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ContentSlot{}, apperr.ErrValidation
		}
		return domain.ContentSlot{}, err
	}
	slot.ID = id
	return slot, nil
}

func (r *Repository) ListContentSlots() []domain.ContentSlot {
	rows, err := r.db.Query(`
SELECT id::text, slot_key, name, COALESCE(description,''), enabled
FROM content_slots
ORDER BY slot_key ASC
`)
	if err != nil {
		return []domain.ContentSlot{}
	}
	defer rows.Close()
	items := make([]domain.ContentSlot, 0)
	for rows.Next() {
		var slot domain.ContentSlot
		if err := rows.Scan(&slot.ID, &slot.SlotKey, &slot.Name, &slot.Description, &slot.Enabled); err == nil {
			items = append(items, slot)
		}
	}
	return items
}

func (r *Repository) CreateSlotItem(slotKey string, item domain.SlotItem) (domain.SlotItem, error) {
	var id string
	err := r.db.QueryRow(`
INSERT INTO content_slot_items (slot_key, content_type, content_id, order_num, enabled)
VALUES ($1,$2,$3,$4,$5)
RETURNING id::text
`, slotKey, item.ContentType, item.ContentID, item.OrderNum, item.Enabled).Scan(&id)
	if err != nil {
		if isForeignKeyViolation(err) {
			return domain.SlotItem{}, apperr.ErrSlotNotFound
		}
		return domain.SlotItem{}, err
	}
	item.ID = id
	item.SlotKey = slotKey
	return item, nil
}

func (r *Repository) ListSlotItems(slotKey string) ([]domain.SlotItem, bool) {
	var exists bool
	if err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM content_slots WHERE slot_key=$1)`, slotKey).Scan(&exists); err != nil || !exists {
		return nil, false
	}
	rows, err := r.db.Query(`
SELECT id::text, slot_key, content_type, content_id::text, order_num, enabled
FROM content_slot_items
WHERE slot_key=$1
ORDER BY order_num ASC, id ASC
`, slotKey)
	if err != nil {
		return []domain.SlotItem{}, true
	}
	defer rows.Close()
	items := make([]domain.SlotItem, 0)
	for rows.Next() {
		var item domain.SlotItem
		if err := rows.Scan(&item.ID, &item.SlotKey, &item.ContentType, &item.ContentID, &item.OrderNum, &item.Enabled); err == nil {
			items = append(items, item)
		}
	}
	return items, true
}

func (r *Repository) DeleteSlotItem(slotKey, itemID string) bool {
	res, err := r.db.Exec(`DELETE FROM content_slot_items WHERE slot_key=$1 AND id=$2`, slotKey, itemID)
	if err != nil {
		return false
	}
	n, _ := res.RowsAffected()
	return n > 0
}

func (r *Repository) ListSlotContent(slotKey string, limit int) ([]domain.SlotContentItem, bool) {
	var exists bool
	if err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM content_slots WHERE slot_key=$1)`, slotKey).Scan(&exists); err != nil || !exists {
		return nil, false
	}
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.db.Query(`
SELECT csi.content_type, csi.content_id::text,
       COALESCE(a.title,''), COALESCE(a.slug,''), COALESCE(a.summary,''),
       COALESCE(m.content,'')
FROM content_slot_items csi
LEFT JOIN articles a ON csi.content_type='article' AND a.id=csi.content_id AND a.status='published' AND a.deleted_at IS NULL
LEFT JOIN moments m ON csi.content_type='moment' AND m.id=csi.content_id AND m.status='published' AND m.deleted_at IS NULL
WHERE csi.slot_key=$1 AND csi.enabled=TRUE
ORDER BY csi.order_num ASC
LIMIT $2
`, slotKey, limit)
	if err != nil {
		return []domain.SlotContentItem{}, true
	}
	defer rows.Close()
	items := make([]domain.SlotContentItem, 0)
	for rows.Next() {
		var it domain.SlotContentItem
		if err := rows.Scan(&it.ContentType, &it.ID, &it.Title, &it.Slug, &it.Summary, &it.Content); err == nil {
			if it.ContentType == "article" && it.Slug == "" {
				continue
			}
			if it.ContentType == "moment" && it.Content == "" {
				continue
			}
			items = append(items, it)
		}
	}
	return items, true
}
