// File: app/stores/site.ts
// Purpose: Global state for Site Settings (Pinia store) to avoid duplicated SSR fetches.

import { defineStore } from 'pinia'
import type { SiteSettings } from '~/composables/useApi'

export const useSiteStore = defineStore('site', {
    state: () => ({
        settings: null as SiteSettings | null,
        navigation: [] as Array<any>,
        socialLinks: [] as Array<any>,
        loaded: false,
    }),
    actions: {
        async fetchAll() {
            if (this.loaded) return

            const api = useApi()
            try {
                const [sets, nav, social] = await Promise.all([
                    api.getSiteSettings(),
                    api.getNav(),
                    api.getSocialLinks(),
                ])
                this.settings = sets
                this.navigation = nav
                this.socialLinks = social
                this.loaded = true
            } catch (err) {
                console.error('Failed to load site foundation data', err)
            }
        }
    }
})
