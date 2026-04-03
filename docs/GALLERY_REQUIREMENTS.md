# Gallery Requirements (v1)

## Purpose
Define the product requirements and interaction rules for the built-in photo gallery module.  
This document focuses on user-facing behavior, admin workflows, image processing, and visibility rules.  
Data schema and API details should be derived from this document in `docs/DATA_MODEL.md` and `docs/API_CONTRACT.md`.

## Product Positioning
- Gallery is a first-class module inside the current blog system, not a separately deployed standalone site.
- Gallery shares the same public domain, top navigation, admin authentication, visual style, SEO baseline, and deployment stack as the blog.
- Gallery content is a global photo stream with tags and per-photo status control. Albums are not included in v1.

## Public Navigation and Routes
- Add a `Gallery` entry to the public top navigation bar.
- Gallery stream page route: `/gallery`.
- Single photo route: `/gallery/{slug}`.
- Public photo URLs use `slug` instead of UUID for readability and sharing.
- If a user opens `/gallery/{slug}` directly, render the standalone photo viewer immediately without list-card transition animation.
- If a user opens a photo from the `/gallery` waterfall stream, animate the clicked card into the viewer with a shared-element style transition, and reverse that transition when closing.

## Gallery Stream Page
- The gallery stream uses a masonry/waterfall layout.
- All photo cards must preserve the original image aspect ratio.
- The gallery stream supports tag-based filtering.
- When no tag filter is active, photo ordering follows the global photo stream order.
- When a tag filter is active, the stream only shows matching photos and viewer previous/next navigation must follow the filtered result order.
- A floating context label in the top-left corner shows the shooting time range and city summary of photos currently visible in the viewport.
  - Time range is calculated from photo shooting time when available.
  - City summary displays city names from currently visible photos.
  - When too many cities are visible, show the first two city names and a `+N` suffix.
- On desktop, hovering a photo card shows only concise overlay information, such as city, camera model, and a small number of tags.
- On mobile, stream cards do not rely on hover-only interactions.

## Photo Viewer
- Clicking a photo opens a full viewer page/state for that photo.
- The large photo must preserve the original aspect ratio and should not be visually cropped by default.
- Desktop viewer layout:
  - large photo on the left/main area
  - detail panel slides in from the right
  - close button in the top-right corner returns to the gallery stream
- Mobile viewer layout:
  - large photo is shown first
  - a `Details` button in the top-left corner opens the metadata panel as a drawer/sheet
  - close button in the top-right corner returns to the gallery stream
- The detail panel can display:
  - title
  - description
  - city/location display name
  - shooting time
  - camera make/model
  - lens model
  - focal length
  - aperture
  - shutter speed
  - ISO
  - tags
- Public photo metadata rendering rule:
  - for each controlled metadata field, the field is shown only when the value exists and the corresponding display switch is enabled
  - hidden fields must be omitted by the public API response, not merely ignored by frontend rendering

## Visual Style
- Gallery UI should follow a frosted-glass visual direction.
- Recommended usage includes:
  - card hover overlays
  - floating top-left context labels
  - viewer side panels
  - mobile metadata drawers
- The style should use translucent backgrounds, blur, soft borders, and restrained contrast, while keeping text legible over photography.
- This visual direction should remain consistent with the existing blog shell instead of introducing a separate design language.

## Photo Status and Visibility
- Each photo has a status.
- Recommended v1 statuses:
  - `draft`: uploaded but not publicly visible
  - `published`: visible in the public gallery and available for article usage
  - `hidden`: not listed in the public gallery stream, but still available for article usage when explicitly referenced
- Public gallery stream only lists `published` photos.
- Single photo public access should respect status rules and must not expose `draft` photos.
- Admin interfaces can read and edit all statuses.

## Per-Photo Metadata Display Control
- Metadata display control is configured per photo and per field.
- Recommended display switches:
  - `taken_at_display`
  - `camera_display`
  - `location_display`
  - `exif_display`
  - `tags_display`
- These switches control public API output and public rendering.
- Admin APIs must still return both raw values and display-switch configuration so operators can edit them.
- If a metadata value is empty, public UI should not render that field even when its display switch is enabled.

## Tags
- Photos support multiple tags.
- Tags should be managed as normalized entities rather than only free-form strings, to avoid duplicates such as `Tokyo`, `tokyo`, and `东京`.
- Admin workflows should allow assigning existing tags and creating new tags while editing photo metadata.
- Public gallery supports filtering by tag.
- Viewer metadata may display tags when `tags_display=true`.

## Location and EXIF Rules
- On upload, the backend should extract a whitelist of EXIF fields from the original file.
- Recommended EXIF whitelist:
  - camera make
  - camera model
  - lens model
  - focal length
  - aperture
  - shutter speed
  - ISO
  - shooting time
  - GPS coordinates for temporary parsing only
- GPS coordinates should be used server-side to reverse geocode a city-level location name and structured location fields.
- After extracting metadata and resolving location, GPS and non-whitelisted EXIF data must be removed from the stored image assets to reduce privacy leakage.
- By default, precise original GPS coordinates should not be stored in the database.
- The public frontend should display city-level or manually edited location text, not raw coordinates.
- If GPS is missing or reverse geocoding fails, admin users must be able to manually edit location fields.
- Shooting time preference:
  - use EXIF shooting time when available
  - fallback to upload/created time when EXIF time is unavailable

## Image Asset Strategy
- Each uploaded photo should produce multiple derived assets:
  - `placeholder`: extremely lightweight blurred placeholder data for progressive loading
  - `display`: optimized image for masonry cards and default inline article rendering
  - `large`: cleaned high-quality image for viewer display
- The original unprocessed file should not be publicly served directly in v1.
- The `large` asset should be generated from the source image after metadata cleaning and quality optimization.
- `display` and `large` assets must preserve the original aspect ratio.
- Storage should use the existing Cloudflare R2 foundation.
- Placeholder implementation can use an LQIP-style tiny blurred image or a hash-based representation such as BlurHash/ThumbHash, as long as the frontend can render a smooth blur-to-sharp transition.

## Progressive Loading Experience
- Gallery cards and article-embedded photos should show a blurred placeholder before the `display` image finishes loading.
- After the actual image is ready, the UI should smoothly transition from placeholder to sharp image.
- The placeholder must be small enough to avoid hurting initial load performance.
- The implementation should avoid layout shift by reserving space according to known width/height or aspect ratio metadata.

## Admin Upload Workflow
- Admin users can batch upload multiple image files and folders.
- Upload workflow should create server-side records/assets first, preferably as `draft`, then allow metadata review and publication from an upload review list.
- After file selection and upload, admin users need a list view showing each uploaded photo and its parsed metadata.
- In that review list, each photo can be previewed and edited individually:
  - title
  - slug
  - description
  - location text
  - shooting time
  - tags
  - status
  - metadata display switches
- The review list should also expose processing state and failures per photo, such as upload success, EXIF parse result, reverse geocoding result, placeholder/display/large generation status, and retry actions.
- Admin users should be able to:
  - remove a photo from the pending batch
  - retry failed processing for a single photo
  - save all as draft
  - publish selected photos in batch
  - batch-apply common metadata changes such as tags, status, or display switches
- Slug conflicts during batch upload/edit should be handled safely by backend validation and either automatic suffixing or explicit admin correction before publication.

## Admin Gallery Management
- Beyond the upload review queue, provide a long-lived gallery management page for all existing photos.
- Required management capabilities:
  - search by title/slug/tag/location
  - filter by status and tag
  - edit single photo metadata and display switches
  - batch status updates
  - batch tag operations
  - reorder or pin photos if manual stream curation is needed
- The upload review queue and the long-lived gallery management page should be separate interfaces to keep batch ingestion and daily content maintenance clear.

## Article Integration
- Article editing should allow inserting an existing gallery photo.
- Article content should reference a stable photo identifier/slug-backed asset model rather than hardcoding a final CDN URL when possible, so future reprocessing or storage migration does not break articles.
- In article rendering:
  - default inline display uses the `display` asset
  - clicking the image opens the `large` asset or the corresponding gallery viewer experience
  - progressive placeholder loading should match the gallery stream behavior
- When article pages display gallery-photo metadata, the same public visibility rule applies: only fields with values and enabled display switches should be exposed.
- `hidden` photos may be usable inside articles if explicitly referenced, but must not appear in the public gallery stream.

## Deletion and Reference Safety
- If a photo is referenced by one or more articles, direct hard delete should be blocked by default.
- Preferred options:
  - prevent deletion until references are removed
  - or use soft deletion with clear admin warnings and safe frontend fallback behavior
- The chosen reference policy must avoid broken article images and should be reflected in API and data-model design.

## Non-Goals for v1
- No album hierarchy in v1.
- No map page or exact coordinate display in v1.
- No public serving of untouched original camera files in v1.
- No separate gallery deployment or separate domain in v1.

## Open Implementation Decisions
- Final placeholder technology: tiny LQIP image, BlurHash, or ThumbHash.
- Reverse geocoding provider and fallback strategy.
- Exact slug auto-generation and conflict-resolution policy.
- Exact public viewer keyboard/swipe navigation behavior.
- Whether photo manual pinning is required in v1 or can be deferred until after the initial gallery launch.
