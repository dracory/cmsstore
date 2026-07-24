const PageMediaApp = {
  data() {
    return {
      loading: true,
      uploading: false,
      deleting: false,
      saving: false,
      isDragOver: false,
      uploadProgress: 0,
      files: [],
      pageId: '',
      draggedIndex: null,
      dragOverIndex: null,
      showAddModal: false,
      showEditModal: false,
      newMediaUrl: '',
      newMediaFileName: '',
      editIndex: null,
      editForm: {
        id: '',
        name: '',
        url: ''
      }
    };
  },

  mounted() {
    if (typeof pageID !== 'undefined') {
      this.pageId = pageID;
    }
    this.loadFiles();
  },

  methods: {
    async loadFiles() {
      this.loading = true;
      try {
        const response = await fetch(urlMediaLoad, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ page_id: this.pageId })
        });
        const data = await response.json();
        if (data.status === 'success') {
          this.files = data.data?.files || [];
        } else {
          Swal.fire({ icon: 'error', title: 'Error', text: data.message || 'Failed to load files' });
        }
      } catch (error) {
        console.error('Error loading files:', error);
        Swal.fire({ icon: 'error', title: 'Error', text: 'Failed to load files' });
      } finally {
        this.loading = false;
      }
    },

    handleFileSelect(event) {
      const files = event.target.files;
      if (files && files.length > 0) {
        this.uploadFiles(files);
      }
      event.target.value = '';
    },

    handleDrop(event) {
      this.isDragOver = false;
      const files = event.dataTransfer.files;
      if (files && files.length > 0) {
        this.uploadFiles(files);
      }
    },

    async addMediaByUrl() {
      if (!this.newMediaUrl.trim()) return;

      const url = this.newMediaUrl.trim();
      const fileName = this.newMediaFileName.trim() || this.getFileNameFromUrl(url);
      const mediaType = this.determineMediaType(url);

      const formData = new FormData();
      formData.append('page_id', this.pageId);
      formData.append('media_url', url);
      formData.append('media_file_name', fileName);
      formData.append('media_type', mediaType);

      this.uploading = true;
      this.uploadProgress = 0;

      try {
        const response = await fetch(urlMediaAdd, {
          method: 'POST',
          body: formData
        });
        const data = await response.json();
        if (data.status === 'success') {
          Swal.fire({ icon: 'success', title: 'Success', text: 'Media added successfully', position: 'top-end', timer: 2000, timerProgressBar: true, showConfirmButton: false });
          this.newMediaUrl = '';
          this.newMediaFileName = '';
          this.showAddModal = false;
          this.loadFiles();
        } else {
          Swal.fire({ icon: 'error', title: 'Error', text: data.message || 'Failed to add media' });
        }
      } catch (error) {
        console.error('Error adding media:', error);
        Swal.fire({ icon: 'error', title: 'Error', text: 'Failed to add media' });
      } finally {
        this.uploading = false;
      }
    },

    determineMediaType(url) {
      if (!url) return 'file';
      const extension = url.split('.').pop().toLowerCase();
      const imageExtensions = ['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg'];
      const videoExtensions = ['mp4', 'webm', 'ogg', 'mov'];
      const audioExtensions = ['mp3', 'wav', 'ogg'];
      if (imageExtensions.includes(extension)) return 'image';
      if (videoExtensions.includes(extension)) return 'video';
      if (audioExtensions.includes(extension)) return 'audio';
      return 'file';
    },

    getFileNameFromUrl(url) {
      if (!url) return '';
      const parts = url.split('/');
      return parts[parts.length - 1] || 'Unnamed';
    },

    async uploadFiles(files) {
      this.uploading = true;
      this.uploadProgress = 0;

      const formData = new FormData();
      formData.append('page_id', this.pageId);
      for (let i = 0; i < files.length; i++) {
        formData.append('files[]', files[i]);
      }

      try {
        const xhr = new XMLHttpRequest();
        xhr.open('POST', urlMediaUpload);

        xhr.upload.onprogress = (e) => {
          if (e.lengthComputable) {
            this.uploadProgress = Math.round((e.loaded / e.total) * 100);
          }
        };

        xhr.onload = () => {
          try {
            const data = JSON.parse(xhr.responseText);
            if (data.status === 'success') {
              Swal.fire({ icon: 'success', title: 'Success', text: 'Files uploaded successfully', position: 'top-end', timer: 3000, timerProgressBar: true, showConfirmButton: false });
              this.loadFiles();
            } else {
              Swal.fire({ icon: 'error', title: 'Error', text: data.message || 'Failed to upload files' });
            }
          } catch (e) {
            Swal.fire({ icon: 'error', title: 'Error', text: 'Failed to parse upload response' });
          }
          this.uploading = false;
          if (this.showAddModal) {
            this.showAddModal = false;
          }
        };

        xhr.onerror = () => {
          Swal.fire({ icon: 'error', title: 'Error', text: 'Failed to upload files' });
          this.uploading = false;
        };

        xhr.send(formData);
      } catch (error) {
        console.error('Error uploading files:', error);
        Swal.fire({ icon: 'error', title: 'Error', text: 'Failed to upload files' });
        this.uploading = false;
      }
    },

    confirmDelete(file) {
      Swal.fire({
        title: 'Delete file?',
        text: file.name,
        icon: 'warning',
        showCancelButton: true,
        confirmButtonText: 'Yes, delete',
        cancelButtonText: 'Cancel'
      }).then((result) => {
        if (result.isConfirmed) {
          this.deleteFile(file);
        }
      });
    },

    async deleteFile(file) {
      this.deleting = true;
      try {
        const response = await fetch(urlMediaDelete, {
          method: 'POST',
          headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
          body: new URLSearchParams({ file_id: file.id })
        });
        const data = await response.json();
        if (data.status === 'success') {
          Swal.fire({ icon: 'success', title: 'Success', text: 'File deleted successfully', position: 'top-end', timer: 3000, timerProgressBar: true, showConfirmButton: false });
          this.loadFiles();
        } else {
          Swal.fire({ icon: 'error', title: 'Error', text: data.message || 'Failed to delete file' });
        }
      } catch (error) {
        console.error('Error deleting file:', error);
        Swal.fire({ icon: 'error', title: 'Error', text: 'Failed to delete file' });
      } finally {
        this.deleting = false;
      }
    },

    async saveMedia() {
      if (this.saving) return;
      this.saving = true;
      try {
        const payload = {
          page_id: this.pageId,
          files: this.files.map((f, i) => ({
            id: f.id,
            name: f.name,
            sequence: i
          }))
        };
        const response = await fetch(urlMediaSave, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        });
        const data = await response.json();
        if (data.status === 'success') {
          Swal.fire({ icon: 'success', title: 'Success', text: 'Media saved successfully', position: 'top-end', timer: 2000, timerProgressBar: true, showConfirmButton: false });
        } else {
          Swal.fire({ icon: 'error', title: 'Error', text: data.message || 'Failed to save media' });
        }
      } catch (error) {
        console.error('Error saving media:', error);
        Swal.fire({ icon: 'error', title: 'Error', text: 'Failed to save media' });
      } finally {
        this.saving = false;
      }
    },

    isImage(file) {
      if (file.type && file.type.startsWith('image/')) return true;
      return this.isImageUrl(file.url || '');
    },

    isImageUrl(url) {
      if (!url) return false;
      if (url.startsWith('data:image/')) return true;
      const ext = url.split('.').pop().toLowerCase();
      return ['jpg', 'jpeg', 'png', 'gif', 'svg', 'webp'].includes(ext);
    },

    dragStart(event, index) {
      this.draggedIndex = index;
      event.dataTransfer.effectAllowed = 'move';
      event.dataTransfer.setData('text/plain', index);
    },

    dragOver(event, index) {
      event.preventDefault();
      event.dataTransfer.dropEffect = 'move';
      if (this.draggedIndex !== null && this.draggedIndex !== index) {
        this.dragOverIndex = index;
      }
    },

    dragEnter(event, index) {
      event.preventDefault();
      if (this.draggedIndex !== null && this.draggedIndex !== index) {
        this.dragOverIndex = index;
      }
    },

    dragLeave(event, index) {
      if (this.dragOverIndex === index) {
        this.dragOverIndex = null;
      }
    },

    drop(event, dropIndex) {
      event.preventDefault();
      if (this.draggedIndex === null || this.draggedIndex === dropIndex) {
        this.draggedIndex = null;
        this.dragOverIndex = null;
        return;
      }
      const draggedItem = this.files[this.draggedIndex];
      this.files.splice(this.draggedIndex, 1);
      this.files.splice(dropIndex, 0, draggedItem);
      this.draggedIndex = null;
      this.dragOverIndex = null;
      this.saveMedia();
    },

    openEditModal(index) {
      this.editIndex = index;
      this.editForm = {
        id: this.files[index].id,
        name: this.files[index].name,
        url: this.files[index].url
      };
      this.showEditModal = true;
    },

    closeEditModal() {
      this.showEditModal = false;
      this.editIndex = null;
      this.editForm = { id: '', name: '', url: '' };
    },

    async saveEdit() {
      if (!this.editForm.name.trim()) return;
      this.files[this.editIndex].name = this.editForm.name.trim();
      this.closeEditModal();
      await this.saveMedia();
    },

    getFileIcon(name) {
      const ext = name.split('.').pop().toLowerCase();
      const iconMap = {
        jpg: 'bi bi-file-image',
        jpeg: 'bi bi-file-image',
        png: 'bi bi-file-image',
        gif: 'bi bi-file-image',
        svg: 'bi bi-file-image',
        webp: 'bi bi-file-image',
        pdf: 'bi bi-file-pdf',
        doc: 'bi bi-file-word',
        docx: 'bi bi-file-word',
        xls: 'bi bi-file-excel',
        xlsx: 'bi bi-file-excel',
        ppt: 'bi bi-file-ppt',
        pptx: 'bi bi-file-ppt',
        mp4: 'bi bi-file-play',
        avi: 'bi bi-file-play',
        mov: 'bi bi-file-play',
        mkv: 'bi bi-file-play',
        mp3: 'bi bi-file-music',
        wav: 'bi bi-file-music',
        flac: 'bi bi-file-music',
        zip: 'bi bi-file-zip',
        rar: 'bi bi-file-zip',
        '7z': 'bi bi-file-zip',
        txt: 'bi bi-file-text',
        md: 'bi bi-file-text',
        csv: 'bi bi-file-spreadsheet',
      };
      return iconMap[ext] || 'bi bi-file-earmark';
    },

    formatSize(sizeStr) {
      const size = parseInt(sizeStr, 10);
      if (isNaN(size) || size === 0) return '';
      const units = ['B', 'KB', 'MB', 'GB'];
      let i = 0;
      let s = size;
      while (s >= 1024 && i < units.length - 1) {
        s /= 1024;
        i++;
      }
      return s.toFixed(i === 0 ? 0 : 1) + ' ' + units[i];
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

function mountMediaApp() {
  const el = document.getElementById('page-media-app');
  const tpl = document.getElementById('page-media-template');
  if (el && tpl) {
    const { createApp } = Vue;
    PageMediaApp.template = tpl.innerHTML;
    createApp(PageMediaApp).mount(el);
  }
}

function initMediaApp() {
  loadSwal(function() {
    loadVue(mountMediaApp);
  });
}

document.addEventListener('DOMContentLoaded', initMediaApp);
