const PageContentApp = {
  data() {
    return {
      loading: true,
      saving: false,
      pageId: '',
      form: {
        title: '',
        content: ''
      }
    };
  },

  mounted() {
    if (typeof pageID !== 'undefined') {
      this.pageId = pageID;
    }
    this.loadContent();
  },

  methods: {
    async loadContent() {
      this.loading = true;
      try {
        const response = await fetch(urlContentLoad, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ page_id: this.pageId })
        });
        const data = await response.json();
        if (data.status === 'success') {
          this.form.title = data.data?.title || '';
          this.form.content = data.data?.content || '';
        } else {
          Swal.fire({ icon: 'error', title: 'Error', text: data.message || 'Failed to load content' });
        }
      } catch (error) {
        console.error('Error loading content:', error);
        Swal.fire({ icon: 'error', title: 'Error', text: 'Failed to load content' });
      } finally {
        this.loading = false;
      }
    },

    async saveContent() {
      if (this.saving) return;
      this.saving = true;
      try {
        const response = await fetch(urlContentSave, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            page_id: this.pageId,
            page_title: this.form.title,
            page_content: this.form.content
          })
        });
        const data = await response.json();
        if (data.status === 'success') {
          Swal.fire({ icon: 'success', title: 'Success', text: 'Page saved successfully', position: 'top-end', timer: 3000, timerProgressBar: true, showConfirmButton: false });
        } else {
          Swal.fire({ icon: 'error', title: 'Error', text: data.message || 'Failed to save content' });
        }
      } catch (error) {
        console.error('Error saving content:', error);
        Swal.fire({ icon: 'error', title: 'Error', text: 'Failed to save content' });
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

function mountContentApp() {
  const el = document.getElementById('page-content-app');
  const tpl = document.getElementById('page-content-template');
  if (el && tpl) {
    const { createApp } = Vue;
    PageContentApp.template = tpl.innerHTML;
    createApp(PageContentApp).mount(el);
  }
}

function initContentApp() {
  loadSwal(function() {
    loadVue(mountContentApp);
  });
}

document.addEventListener('DOMContentLoaded', initContentApp);
