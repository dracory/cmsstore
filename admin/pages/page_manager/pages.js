const { createApp } = Vue;

/**
 * CmsPagesApp is a Vue.js component for managing CMS pages.
 * It provides a table view with filtering, sorting, pagination,
 * and CRUD operations for pages.
 */
const CmsPagesApp = {
  data() {
    return {
      loading: true,
      showCreateModal: false,
      showFilterModal: false,
      creating: false,

      pages: [],
      totalPages: 0,
      sites: [],

      currentPage: 0,
      perPage: 20,

      filters: {
        search: '',
        status: '',
        siteId: '',
        dateFrom: '',
        dateTo: ''
      },

      sortByColumn: 'created_at',
      sortOrder: 'desc',

      createForm: {
        name: '',
        siteId: ''
      }
    };
  },

  computed: {
    visiblePages() {
      const pages = [];
      const total = Math.ceil(this.totalPages / this.perPage);
      const start = Math.max(0, this.currentPage - 2);
      const end = Math.min(total - 1, this.currentPage + 2);

      for (let i = start; i <= end; i++) {
        pages.push(i);
      }
      return pages;
    },

    totalPagesCount() {
      return Math.ceil(this.totalPages / this.perPage);
    },

    filterStatus() {
      const parts = [];
      if (this.filters.search) parts.push(`search: "${this.filters.search}"`);
      if (this.filters.status) parts.push(`status: ${this.filters.status}`);
      if (this.filters.dateFrom) parts.push(`from: ${this.filters.dateFrom}`);
      if (this.filters.dateTo) parts.push(`to: ${this.filters.dateTo}`);

      if (parts.length === 0) return 'Showing all pages';
      return 'Showing pages with ' + parts.join(', ');
    }
  },

  mounted() {
    const urlParams = new URLSearchParams(window.location.search);

    this.filters.search = urlParams.get('search') || '';
    this.filters.status = urlParams.get('status') || '';
    this.filters.siteId = urlParams.get('site_id') || '';
    this.filters.dateFrom = urlParams.get('date_from') || '';
    this.filters.dateTo = urlParams.get('date_to') || '';
    this.sortByColumn = urlParams.get('sort_by') || 'created_at';
    this.sortOrder = urlParams.get('sort_order') || 'desc';
    this.currentPage = parseInt(urlParams.get('page') || '0', 10);
    this.perPage = parseInt(urlParams.get('per_page') || '20', 10);

    this.loadPages();
  },

  methods: {
    async loadPages() {
      this.loading = true;
      try {
        const response = await fetch(urlPagesLoad, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            page: this.currentPage,
            per_page: this.perPage,
            search: this.filters.search,
            status: this.filters.status,
            site_id: this.filters.siteId,
            date_from: this.filters.dateFrom,
            date_to: this.filters.dateTo,
            sort_order: this.sortOrder,
            sort_by: this.sortByColumn
          })
        });
        const data = await response.json();

        if (data.status === 'success') {
          this.pages = data.data?.pages || [];
          this.totalPages = data.data?.total || 0;
          this.sites = data.data?.sites || [];
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to load pages'
          });
        }
      } catch (error) {
        console.error('Error loading pages:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: error.message || 'Failed to load pages'
        });
      } finally {
        this.loading = false;
      }
    },

    sortBy(column) {
      if (this.sortByColumn === column) {
        this.sortOrder = this.sortOrder === 'asc' ? 'desc' : 'asc';
      } else {
        this.sortByColumn = column;
        this.sortOrder = 'asc';
      }
      this.currentPage = 0;
      this.applyFilters();
    },

    goToPage(page) {
      const total = Math.ceil(this.totalPages / this.perPage);
      if (page < 0 || page >= total) return;
      this.currentPage = page;
      this.applyFilters();
    },

    async deletePage(page) {
      const result = await Swal.fire({
        icon: 'warning',
        title: 'Delete Page?',
        text: `Are you sure you want to delete "${page.name}"?`,
        showCancelButton: true,
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        confirmButtonColor: '#dc3545'
      });

      if (!result.isConfirmed) return;

      try {
        const response = await fetch(urlPageDelete, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ page_id: page.id })
        });

        const data = await response.json();

        if (data.status === 'success') {
          Swal.fire({
            icon: 'success',
            title: 'Deleted',
            text: 'Page deleted successfully',
            timer: 1500,
            showConfirmButton: false
          });
          this.loadPages();
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to delete page'
          });
        }
      } catch (error) {
        console.error('Error deleting page:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to delete page'
        });
      }
    },

    formatDate(dateString) {
      if (!dateString) return '-';
      const date = new Date(dateString);
      const options = { day: 'numeric', month: 'short', year: 'numeric' };
      return date.toLocaleDateString('en-GB', options);
    },

    getPageUpdateUrl(pageId) {
      return urlPageUpdate.replace('PAGE_ID_PLACEHOLDER', pageId);
    },

    openCreateModal() {
      this.createForm.name = '';
      this.createForm.siteId = '';
      this.showCreateModal = true;
    },

    closeCreateModal() {
      this.showCreateModal = false;
      this.createForm.name = '';
      this.createForm.siteId = '';
    },

    openFilterModal() {
      this.showFilterModal = true;
    },

    closeFilterModal() {
      this.showFilterModal = false;
    },

    applyFilters() {
      const urlParams = new URLSearchParams(window.location.search);
      const path = urlParams.get('path');

      const params = new URLSearchParams();
      if (path) params.set('path', path);
      if (this.filters.search) params.set('search', this.filters.search);
      if (this.filters.status) params.set('status', this.filters.status);
      if (this.filters.siteId) params.set('site_id', this.filters.siteId);
      if (this.filters.dateFrom) params.set('date_from', this.filters.dateFrom);
      if (this.filters.dateTo) params.set('date_to', this.filters.dateTo);
      params.set('page', 0);
      params.set('per_page', this.perPage);
      params.set('sort_order', this.sortOrder);
      params.set('sort_by', this.sortByColumn);

      const newUrl = `${window.location.pathname}?${params.toString()}`;
      window.history.pushState({}, '', newUrl);

      this.closeFilterModal();
      this.loadPages();
    },

    clearFilters() {
      this.filters = {
        search: '',
        status: '',
        siteId: '',
        dateFrom: '',
        dateTo: ''
      };
      this.applyFilters();
    },

    async createPage() {
      if (!this.createForm.name || !this.createForm.siteId) return;

      this.creating = true;
      try {
        const response = await fetch(urlPageCreate, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            name: this.createForm.name,
            site_id: this.createForm.siteId
          })
        });

        const data = await response.json();

        if (data.status === 'success') {
          Swal.fire({
            icon: 'success',
            title: 'Success',
            text: 'Page created successfully',
            timer: 1500,
            showConfirmButton: false
          });
          this.closeCreateModal();
          window.open(urlPageUpdate.replace('PAGE_ID_PLACEHOLDER', data.data.id), '_blank');
          this.loadPages();
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to create page'
          });
        }
      } catch (error) {
        console.error('Error creating page:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to create page'
        });
      } finally {
        this.creating = false;
      }
    }
  }
};

document.addEventListener('DOMContentLoaded', () => {
  const el = document.getElementById('cms-pages-app');
  if (el) {
    createApp(CmsPagesApp).mount('#cms-pages-app');
  }
});
