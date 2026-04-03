/**
 * Purpose: Define gallery domain types for admin gallery management module.
 * Module: frontend-admin-react/types, gallery domain types layer.
 * Related: gallery API module, GalleryPage, and GalleryEditorPage.
 */

export type PhotoStatus = 'draft' | 'published' | 'hidden';
export type ProcessingStatus = 'pending' | 'processing' | 'completed' | 'failed';

export interface GalleryPhoto {
  id: string;
  title: string;
  slug: string;
  description: string;
  status: PhotoStatus;
  locationName: string;
  locationCity: string;
  locationState: string;
  locationCountry: string;
  takenAt?: string;
  cameraMake: string;
  cameraModel: string;
  lensModel: string;
  focalLength: string;
  aperture: string;
  shutterSpeed: string;
  iso: string;
  fileSizeBytes: number;
  width: number;
  height: number;
  takenAtDisplay: boolean;
  cameraDisplay: boolean;
  locationDisplay: boolean;
  exifDisplay: boolean;
  tagsDisplay: boolean;
  placeholderData: string;
  displayUrl: string;
  largeUrl: string;
  processingStatus: ProcessingStatus;
  processingError: string;
  sortOrder: number;
  articleRefCount: number;
  tagSlugs: string[];
  createdAt: string;
  updatedAt: string;
}

export interface GalleryTag {
  id: string;
  name: string;
  slug: string;
}

export interface GalleryPhotoListParams {
  page?: number;
  pageSize?: number;
  status?: PhotoStatus | '';
  tag?: string;
  keyword?: string;
}

export interface GalleryPhotoFormValues {
  title: string;
  slug: string;
  description?: string;
  status: PhotoStatus;
  locationName?: string;
  locationCity?: string;
  locationState?: string;
  locationCountry?: string;
  takenAt?: string;
  cameraMake?: string;
  cameraModel?: string;
  lensModel?: string;
  focalLength?: string;
  aperture?: string;
  shutterSpeed?: string;
  iso?: string;
  takenAtDisplay?: boolean;
  cameraDisplay?: boolean;
  locationDisplay?: boolean;
  exifDisplay?: boolean;
  tagsDisplay?: boolean;
  sortOrder?: number;
  tagSlugs?: string[];
}
