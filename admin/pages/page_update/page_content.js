const PageContentApp = {
  data() {
    return {
      loading: true,
      saving: false,
      pageId: '',
      editor: '',
      codeMirrorInstance: null,
      trumbowygInstance: null,
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
          this.editor = data.data?.editor || '';
        } else {
          Swal.fire({ icon: 'error', title: 'Error', text: data.message || 'Failed to load content' });
        }
      } catch (error) {
        console.error('Error loading content:', error);
        Swal.fire({ icon: 'error', title: 'Error', text: 'Failed to load content' });
      } finally {
        this.loading = false;
        this.$nextTick(() => {
          if (this.editor === 'codemirror') {
            this.initCodeMirror();
          } else if (this.editor === 'htmlarea') {
            this.initTrumbowyg();
          }
        });
      }
    },

    loadCodeMirrorAssets(callback) {
      if (typeof CodeMirror !== 'undefined') {
        callback();
        return;
      }

      const css = document.createElement('link');
      css.rel = 'stylesheet';
      css.href = '//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/codemirror.min.css';
      document.head.appendChild(css);

      const scripts = [
        '//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/codemirror.min.js',
        '//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/xml/xml.min.js',
        '//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/htmlmixed/htmlmixed.min.js',
        '//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/javascript/javascript.js',
        '//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/css/css.js',
        '//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/clike/clike.min.js',
        '//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/php/php.min.js',
        '//cdnjs.cloudflare.com/ajax/libs/codemirror/2.36.0/formatting.min.js',
        '//cdnjs.cloudflare.com/ajax/libs/codemirror/3.22.0/addon/edit/matchbrackets.min.js'
      ];

      let loaded = 0;
      scripts.forEach((src) => {
        const s = document.createElement('script');
        s.src = src;
        s.onload = () => {
          loaded++;
          if (loaded === scripts.length) {
            callback();
          }
        };
        s.onerror = () => console.error('Failed to load CodeMirror script:', src);
        document.head.appendChild(s);
      });
    },

    loadScript(src, onError) {
      return new Promise((resolve, reject) => {
        const s = document.createElement('script');
        s.src = src;
        s.onload = () => resolve();
        s.onerror = () => { console.error('Failed to load script:', src); reject(new Error('Failed to load: ' + src)); };
        document.head.appendChild(s);
      });
    },

    loadJQuery(callback) {
      if (typeof jQuery !== 'undefined') {
        callback();
        return;
      }
      this.loadScript('https://code.jquery.com/jquery-3.7.1.min.js').then(callback).catch(() => {});
    },

    loadTrumbowygAssets(callback) {
      if (typeof jQuery !== 'undefined' && jQuery.fn && jQuery.fn.trumbowyg) {
        callback();
        return;
      }

      const css = document.createElement('link');
      css.rel = 'stylesheet';
      css.href = 'https://cdnjs.cloudflare.com/ajax/libs/Trumbowyg/2.27.3/ui/trumbowyg.min.css';
      document.head.appendChild(css);

      this.loadJQuery(() => {
        this.loadScript('https://cdnjs.cloudflare.com/ajax/libs/Trumbowyg/2.27.3/trumbowyg.min.js').then(callback).catch(() => {});
      });
    },

    initTrumbowyg() {
      this.loadTrumbowygAssets(() => {
        const ta = document.querySelector('textarea[name="page_content"]');
        if (!ta) return;

        this.trumbowygInstance = jQuery(ta).trumbowyg({
          btns: [
            ['viewHTML'],
            ['undo', 'redo'],
            ['formatting'],
            ['strong', 'em', 'del'],
            ['superscript', 'subscript'],
            ['link', 'justifyLeft', 'justifyRight', 'justifyCenter', 'justifyFull'],
            ['unorderedList', 'orderedList'],
            ['insertImage'],
            ['removeformat'],
            ['horizontalRule'],
            ['fullscreen']
          ],
          autogrow: true,
          removeformatPasted: true,
          tagsToRemove: ['script', 'link', 'embed', 'iframe', 'input'],
          tagsToKeep: ['hr', 'img', 'i'],
          autogrowOnEnter: true,
          linkTargets: ['_blank']
        });

        jQuery(ta).on('tbwchange tbwblur', () => {
          this.form.content = jQuery(ta).trumbowyg('html');
        });
      });
    },

    initCodeMirror() {
      this.loadCodeMirrorAssets(() => {
        const ta = document.querySelector('textarea[name="page_content"]');
        if (!ta) return;

        this.codeMirrorInstance = CodeMirror.fromTextArea(ta, {
          lineNumbers: true,
          matchBrackets: true,
          mode: 'application/x-httpd-php',
          indentUnit: 4,
          indentWithTabs: true,
          enterMode: 'keep',
          tabMode: 'shift'
        });

        this.codeMirrorInstance.on('change', () => {
          this.form.content = this.codeMirrorInstance.getValue();
        });
      });
    },

    async saveContent() {
      if (this.saving) return;
      this.saving = true;

      if (this.codeMirrorInstance) {
        this.form.content = this.codeMirrorInstance.getValue();
      }

      if (this.trumbowygInstance) {
        this.form.content = this.trumbowygInstance.trumbowyg('html');
      }

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
