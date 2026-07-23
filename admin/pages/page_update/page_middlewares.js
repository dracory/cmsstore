const PageMiddlewaresApp = {
  data() {
    return {
      loading: true,
      saving: false,
      pageId: '',
      newBefore: '',
      newAfter: '',
      form: {
        before: [],
        after: []
      }
    };
  },

  mounted() {
    if (typeof pageID !== 'undefined') {
      this.pageId = pageID;
    }
    this.loadMiddlewares();
  },

  methods: {
    async loadMiddlewares() {
      this.loading = true;
      try {
        const response = await fetch(urlMiddlewaresLoad, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ page_id: this.pageId })
        });
        const data = await response.json();
        if (data.status === 'success') {
          this.form.before = data.data?.before || [];
          this.form.after = data.data?.after || [];
        } else {
          Swal.fire({ icon: 'error', title: 'Error', text: data.message || 'Failed to load middlewares' });
        }
      } catch (error) {
        console.error('Error loading middlewares:', error);
        Swal.fire({ icon: 'error', title: 'Error', text: 'Failed to load middlewares' });
      } finally {
        this.loading = false;
      }
    },

    async saveMiddlewares() {
      if (this.saving) return;
      this.saving = true;
      try {
        const response = await fetch(urlMiddlewaresSave, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            page_id: this.pageId,
            middlewares_before: this.form.before,
            middlewares_after: this.form.after
          })
        });
        const data = await response.json();
        if (data.status === 'success') {
          Swal.fire({ icon: 'success', title: 'Success', text: 'Middlewares saved successfully', position: 'top-end', timer: 3000, timerProgressBar: true, showConfirmButton: false });
        } else {
          Swal.fire({ icon: 'error', title: 'Error', text: data.message || 'Failed to save middlewares' });
        }
      } catch (error) {
        console.error('Error saving middlewares:', error);
        Swal.fire({ icon: 'error', title: 'Error', text: 'Failed to save middlewares' });
      } finally {
        this.saving = false;
      }
    },

    addBefore() {
      const mw = this.newBefore.trim();
      if (!mw) return;
      this.form.before.push(mw);
      this.newBefore = '';
    },

    removeBefore(index) {
      this.form.before.splice(index, 1);
    },

    moveUpBefore(index) {
      if (index === 0) return;
      const tmp = this.form.before[index];
      this.form.before[index] = this.form.before[index - 1];
      this.form.before[index - 1] = tmp;
    },

    moveDownBefore(index) {
      if (index === this.form.before.length - 1) return;
      const tmp = this.form.before[index];
      this.form.before[index] = this.form.before[index + 1];
      this.form.before[index + 1] = tmp;
    },

    addAfter() {
      const mw = this.newAfter.trim();
      if (!mw) return;
      this.form.after.push(mw);
      this.newAfter = '';
    },

    removeAfter(index) {
      this.form.after.splice(index, 1);
    },

    moveUpAfter(index) {
      if (index === 0) return;
      const tmp = this.form.after[index];
      this.form.after[index] = this.form.after[index - 1];
      this.form.after[index - 1] = tmp;
    },

    moveDownAfter(index) {
      if (index === this.form.after.length - 1) return;
      const tmp = this.form.after[index];
      this.form.after[index] = this.form.after[index + 1];
      this.form.after[index + 1] = tmp;
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

function mountMiddlewaresApp() {
  const el = document.getElementById('page-middlewares-app');
  const tpl = document.getElementById('page-middlewares-template');
  if (el && tpl) {
    const { createApp } = Vue;
    PageMiddlewaresApp.template = tpl.innerHTML;
    createApp(PageMiddlewaresApp).mount(el);
  }
}

function initMiddlewaresApp() {
  loadSwal(function() {
    loadVue(mountMiddlewaresApp);
  });
}

document.addEventListener('DOMContentLoaded', initMiddlewaresApp);
