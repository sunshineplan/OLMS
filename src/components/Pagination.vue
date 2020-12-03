<template>
  <a class="text-secondary" style="padding-left: 15px" v-if="total">
    <span v-if="total == 1">{{ $t("Showing") }} 1 {{ $t("row") }}</span>
    <span v-else-if="page == current">
      {{ $t("Showing") }} {{ (current - 1) * 10 + 1 }}-{{ total }}
      {{ $t("of") }} {{ total }} {{ $t("rows") }}
    </span>
    <span v-else>
      {{ $t("Showing") }} {{ (current - 1) * 10 + 1 }}-{{ current * 10 }}
      {{ $t("of") }} {{ total }} {{ $t("rows") }}
    </span>
  </a>
  <slot></slot>
  <nav v-if="total">
    <ul class="pagination justify-content-center">
      <li class="page-item" :disabled="current > 1">
        <a class="page-link" :click="current--">{{ $t("Previous") }}</a>
      </li>
      <div v-for="(i, index) in items" :key="i">
        <li class="page-item" v-if="index > 1 && i - items[index - 1] > 1">
          <a class="page-link">...</a>
        </li>
        <li class="page-item" :class="{ active: i == current }">
          <a class="page-link" :click="(current = i)">{{ i }}</a>
        </li>
      </div>
      <li class="page-item" :disabled="current < page">
        <a class="page-link" :click="current++">{{ $t("Next") }}</a>
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
    current: {
      get() {
        return this.$store.state.current;
      },
      set(value) {
        this.$store.commit("current", value);
      },
    },
    page() {
      return Math.ceil(this.total / 10);
    },
    items() {
      return Array.from(
        new Set(
          [1, 2, this.page, this.page - 1]
            .concat(new Array(5).fill().map((d, i) => i + this.current - 2))
            .sort((a, b) => {
              return a - b;
            })
        )
      );
    },
  },
};
</script>