const PageSettingsApp = {
  data() {
    return {
      loading: true,
      saving: false,
      pageId: '',
      sites: [],
      templates: [],
      form: {
        status: '',
        templateId: '',
        editor: '',
        name: '',
        siteId: '',
        memo: ''
      }
    };
  },

  mounted() {
    if (typeof pageID !== 'undefined') {
      this.pageId = pageID;
    }
    this.loadSettings();
  },

  methods: {
    async loadSettings() {
      this.loading = true;
      try {
        const response = await fetch(urlSettingsLoad, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ page_id: this.pageId })
        });
        const data = await response.json();
        if (data.status === 'success') {
          this.form.status = data.data?.status || '';
          this.form.templateId = data.data?.template_id || '';
          this.form.editor = data.data?.editor || '';
          this.form.name = data.data?.name || '';
          this.form.siteId = data.data?.site_id || '';
          this.form.memo = data.data?.memo || '';
          this.sites = data.data?.sites || [];
          this.templates = data.data?.templates || [];
        } else {
          Swal.fire({ icon: 'error', title: 'Error', text: data.message || 'Failed to load settings' });
        }
      } catch (error) {
        console.error('Error loading settings:', error);
        Swal.fire({ icon: 'error', title: 'Error', text: 'Failed to load settings' });
      } finally {
        this.loading = false;
      }
    },

    async saveSettings() {
      if (this.saving) return;
      if (!this.form.status) {
        Swal.fire({ icon: 'error', title: 'Error', text: 'Status is required' });
        return;
      }
      this.saving = true;
      try {
        const response = await fetch(urlSettingsSave, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            page_id: this.pageId,
            page_status: this.form.status,
            page_template_id: this.form.templateId,
            page_editor: this.form.editor,
            page_name: this.form.name,
            page_site_id: this.form.siteId,
            page_memo: this.form.memo
          })
        });
        const data = await response.json();
        if (data.status === 'success') {
          Swal.fire({ icon: 'success', title: 'Success', text: 'Page saved successfully', position: 'top-end', timer: 3000, timerProgressBar: true, showConfirmButton: false });
        } else {
          Swal.fire({ icon: 'error', title: 'Error', text: data.message || 'Failed to save settings' });
        }
      } catch (error) {
        console.error('Error saving settings:', error);
        Swal.fire({ icon: 'error', title: 'Error', text: 'Failed to save settings' });
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

function mountSettingsApp() {
  const el = document.getElementById('page-settings-app');
  const tpl = document.getElementById('page-settings-template');
  if (el && tpl) {
    const { createApp } = Vue;
    PageSettingsApp.template = tpl.innerHTML;
    createApp(PageSettingsApp).mount(el);
  }
}

function initSettingsApp() {
  loadSwal(function() {
    loadVue(mountSettingsApp);
  });
}

document.addEventListener('DOMContentLoaded', initSettingsApp);
