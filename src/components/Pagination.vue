<template>
  <a class="text-secondary" style="padding-left: 15px" v-if="total">
    <span v-if="total == 1">{{ $t("Showing") }} 1 {{ $t("row") }}</span>
    <span v-else-if="page == page">
      {{ $t("Showing") }} {{ (page - 1) * 10 + 1 }}-{{ total }}
      {{ $t("of") }} {{ total }} {{ $t("rows") }}
    </span>
    <span v-else>
      {{ $t("Showing") }} {{ (page - 1) * 10 + 1 }}-{{ page * 10 }}
      {{ $t("of") }} {{ total }} {{ $t("rows") }}
    </span>
  </a>
  <slot></slot>
  <nav v-if="total">
    <ul class="pagination justify-content-center">
      <li class="page-item" :disabled="page > 1">
        <a class="page-link" @click="page--">{{ $t("Previous") }}</a>
      </li>
      <div v-for="(i, index) in items" :key="i">
        <li class="page-item" v-if="index > 1 && i - items[index - 1] > 1">
          <a class="page-link">...</a>
        </li>
        <li class="page-item" :class="{ active: i == page }">
          <a class="page-link" @click="(page = i)">{{ i }}</a>
        </li>
      </div>
      <li class="page-item" :disabled="page < page">
        <a class="page-link" @click="page++">{{ $t("Next") }}</a>
      </li>
    </ul>
  </nav>
</template>

<script>
export default {
  name: "Pagination",
  props: {
    total: Number,
  },
  computed: {
    page: {
      get() {
        return this.$store.state.page;
      },
      set(page) {
        this.$store.commit("page", page);
      },
    },
    page() {
      return Math.ceil(this.total / 10);
    },
    items() {
      return Array.from(
        new Set(
          [1, 2, this.page, this.page - 1]
            .concat(new Array(5).fill().map((d, i) => i + this.page - 2))
            .sort((a, b) => {
              return a - b;
            })
        )
      );
    },
  },
};
</script>