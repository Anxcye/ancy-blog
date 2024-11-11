import { useBrowserStore } from '@/stores/browser'

import type { ComputedRef } from 'vue'

export default function getMeta(
  description: ComputedRef<string>,
  keywords: ComputedRef<string[]>,
  ogTitle?: ComputedRef<string>,
  ogDescription?: ComputedRef<string>,
) {
  return [
    {
      name: 'description',
      content: description,
    },
    {
      name: 'keywords',
      content: [['AnCy', 'Anxcye'], keywords.value].join(','),
    },
    {
      property: 'og:title',
      content: ogTitle ?? useBrowserStore().getTitle,
    },
    {
      property: 'og:description',
      content: ogDescription ?? description,
    },
  ]
}
