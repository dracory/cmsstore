const PageSEOApp = {
  data() {
    return {
      loading: true,
      saving: false,
      pageId: '',
      form: {
        alias: '',
        canonicalUrl: '',
        metaDescription: '',
        metaKeywords: '',
        metaRobots: ''
      }
    };
  },

  mounted() {
    if (typeof pageID !== 'undefined') {
      this.pageId = pageID;
    }
    this.loadSEO();
  },

  methods: {
    async loadSEO() {
      this.loading = true;
      try {
        const response = await fetch(urlSEOLoad, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ page_id: this.pageId })
        });
        const data = await response.json();
        if (data.status === 'success') {
          this.form.alias = data.data?.alias || '';
          this.form.canonicalUrl = data.data?.canonical_url || '';
          this.form.metaDescription = data.data?.meta_description || '';
          this.form.metaKeywords = data.data?.meta_keywords || '';
          this.form.metaRobots = data.data?.meta_robots || '';
        } else {
          Swal.fire({ icon: 'error', title: 'Error', text: data.message || 'Failed to load SEO data' });
        }
      } catch (error) {
        console.error('Error loading SEO:', error);
        Swal.fire({ icon: 'error', title: 'Error', text: 'Failed to load SEO data' });
      } finally {
        this.loading = false;
      }
    },

    async saveSEO() {
      if (this.saving) return;
      this.saving = true;
      try {
        const response = await fetch(urlSEOSave, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            page_id: this.pageId,
            page_alias: this.form.alias,
            page_canonical_url: this.form.canonicalUrl,
            page_meta_description: this.form.metaDescription,
            page_meta_keywords: this.form.metaKeywords,
            page_meta_robots: this.form.metaRobots
          })
        });
        const data = await response.json();
        if (data.status === 'success') {
          Swal.fire({ icon: 'success', title: 'Success', text: 'Page saved successfully', position: 'top-end', timer: 3000, timerProgressBar: true, showConfirmButton: false });
        } else {
          Swal.fire({ icon: 'error', title: 'Error', text: data.message || 'Failed to save SEO data' });
        }
      } catch (error) {
        console.error('Error saving SEO:', error);
        Swal.fire({ icon: 'error', title: 'Error', text: 'Failed to save SEO data' });
      } finally {
        this.saving = false;
      }
    }
  }
};

function loadVue(callback) {
  if (typeof Vue !== 'undefined') {
    callback();
    return;
  }

  const script = document.createElement('script');
  script.src = 'https://unpkg.com/vue@3/dist/vue.global.js';
  script.onload = () => callback();
  script.onerror = () => console.error('Failed to load Vue');
  document.head.appendChild(script);
}

function loadSwal(callback) {
  if (typeof Swal !== 'undefined') {
    callback();
    return;
  }

  const script = document.createElement('script');
  script.src = 'https://cdn.jsdelivr.net/npm/sweetalert2@11';
  script.onload = () => callback();
  script.onerror = () => console.error('Failed to load SweetAlert2');
  document.head.appendChild(script);
}

function mountSEOApp() {
  const el = document.getElementById('page-seo-app');
  const tpl = document.getElementById('page-seo-template');
  if (el && tpl) {
    const { createApp } = Vue;
    PageSEOApp.template = tpl.innerHTML;
    createApp(PageSEOApp).mount(el);
  }
}

function initSEOApp() {
  loadSwal(function() {
    loadVue(mountSEOApp);
  });
}

document.addEventListener('DOMContentLoaded', initSEOApp);
