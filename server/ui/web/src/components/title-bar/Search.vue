<template>
  <div class="search">
    <!--suppress HtmlFormInputWithoutLabel -->
    <input
      class="search__input"
      type="text"
      :placeholder="$t('ui.components.search.placeholder')"
      v-model.trim="searched"
      @keyup.enter="search"
    />
  </div>
</template>

<script lang="ts">
import Component from "vue-class-component";
import Vue from "vue";

@Component
export default class Search extends Vue {
  get searched(): string {
    return this.$store.state.search.searched;
  }

  set searched(searched: string) {
    this.$store.commit("search/setSearched", searched);
  }

  search(): void {
    this.$store.dispatch("search/search", this.$router);
  }
}
</script>

<style lang="scss">
.search {
  &__input {
    display: inline-block;
    width: 100%;
    font-size: 2rem;
    padding: 1.4rem;
    border: none;
    transition: all 0.2s;

    &:hover {
      color: $color-complementary;
      &::placeholder {
        color: $color-complementary;
      }
    }

    &:focus {
      color: $color-complementary-dark;
      &::placeholder {
        color: $color-complementary-dark;
      }
    }
  }
}
</style>
