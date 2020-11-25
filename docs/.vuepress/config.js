module.exports = {
  locales: {
    '/': {
      lang: 'en-US',
      title: 'kuberik',
      description: 'Kubernetes pipeline engine'
    },
  },
  head: [
    ['link', { rel: 'icon', href: '/assets/img/logo.svg' }],
  ],
  themeConfig: {
    logo: '/assets/img/logo.svg',
    nav: [
      { text: 'Usage', link: '/usage/' },
      { text: 'Examples', link: '/examples/' },
      { text: 'Extending', link: '/extending/' },
      { text: 'Contributing', link: '/contributing/' },
    ],
    sidebar: {
      '/usage/': getUsageSidebar(),
      '/extending/': getExtendingSidebar(),
      '/contributing/': getContributingSidebar(),
    },
    searchPlaceholder: 'Search...',
    // Assumes GitHub. Can also be a full GitLab url.
    repo: 'kuberik/engine',
    // Customising the header label
    // Defaults to "GitHub"/"GitLab"/"Bitbucket" depending on `themeConfig.repo`
    repoLabel: 'Contribute!',

    // if your docs are not at the root of the repo:
    docsDir: 'docs',
    // if your docs are in a specific branch (defaults to 'master'):
    docsBranch: 'master',
    // defaults to false, set to true to enable
    editLinks: true,
    // default value is true. Allows to hide next page links on all pages
    nextLinks: false,
    // default value is true. Allows to hide prev page links on all pages
    prevLinks: false,
    // custom text for edit link. Defaults to "Edit this page"
    editLinkText: 'Help us improve this page!'
  }
}

function getContributingSidebar() {
  return [{
    title: "Contributing",
    collapsable: false,
    sidebarDepth: 2,
    children: [
      'core-principles',
      'architecture',
      'design-goals',
      'non-goals',
      'pipeline-features',
    ]
  }]
}

function getUsageSidebar() {
  return [{
    title: "Usage",
    collapsable: false,
    sidebarDepth: 2,
    children: [
      ['', 'Introduction'],
      'terminology',
      'screenplay-reference',
      'screeners',
    ]
  }, {
    title: "Advanced",
    collapsable: false,
    sidebarDepth: 2,
    children: [
      'api-reference',
    ]
  }]
}

function getExtendingSidebar() {
  return [{
    title: "Extending",
    collapsable: false,
    sidebarDepth: 2,
    children: [
      'writing-screeners',
    ]
  }]
}
