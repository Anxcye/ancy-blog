// File: styles/naive-theme.ts
// Purpose: Define Naive UI theme overrides for enterprise-level admin visual consistency.
// Module: frontend-admin/styles, design token integration layer.
// Related: App root provider, shell layout, form/table/card components.
import type { GlobalThemeOverrides } from 'naive-ui';

export const lightThemeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#26a69a',
    primaryColorHover: '#219a8f',
    primaryColorPressed: '#1e8b81',
    primaryColorSuppl: '#26a69a',
    borderColor: '#e4ece9',
    borderRadius: '12px',
    borderRadiusSmall: '10px',
    textColorBase: '#1b2631',
    bodyColor: '#f4f7f6',
    cardColor: '#ffffff',
    modalColor: '#ffffff',
    popoverColor: '#ffffff',
  },
  Card: {
    borderRadius: '14px',
    color: '#ffffff',
  },
  Layout: {
    color: '#f4f7f6',
    siderColor: '#ffffff',
    headerColor: '#ffffff',
  },
  Menu: {
    itemTextColor: '#5f6c76',
    itemTextColorActive: '#0e756b',
    itemTextColorActiveHover: '#0e756b',
    itemColorActive: '#e7f7f4',
    itemColorActiveHover: '#e7f7f4',
  },
};

export const darkThemeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#5ac8be',
    primaryColorHover: '#6fd2c9',
    primaryColorPressed: '#49b6ac',
    primaryColorSuppl: '#5ac8be',
    bodyColor: '#1e252c',
    cardColor: '#262f37',
    modalColor: '#262f37',
    popoverColor: '#262f37',
    tableColor: '#262f37',
    inputColor: '#313b45',
    borderColor: '#3a4551',
    textColorBase: '#e6edf3',
    borderRadius: '12px',
    borderRadiusSmall: '10px',
  },
  Card: {
    borderRadius: '14px',
  },
  Layout: {
    color: '#1e252c',
    siderColor: '#242d36',
    headerColor: '#242d36',
  },
  Input: {
    color: '#313b45',
  },
  Select: {
    peers: {
      InternalSelection: {
        color: '#313b45',
      },
    },
  },
  Menu: {
    itemTextColor: '#a5b2bf',
    itemTextColorActive: '#7bd4cd',
    itemTextColorActiveHover: '#7bd4cd',
    itemColorActive: 'rgba(122, 212, 205, 0.14)',
    itemColorActiveHover: 'rgba(122, 212, 205, 0.14)',
  },
};
