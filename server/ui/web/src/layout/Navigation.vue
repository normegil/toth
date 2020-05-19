<template>
  <div class="navigation">
    <div class="navigation__logo-container">
      <Logo class="navigation__logo" />
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Logo from "../components/navigation/Logo.vue";

@Component({
  components: { Logo }
})
export default class Navigation extends Vue {}
</script>

<style lang="scss">
.navigation {
  background-color: $color-primary-darker;
  position: relative;

  &__logo-container {
    background-color: $color-primary;
    display: flex;
    align-items: baseline;
    padding: 0.6rem;
  }

  &__logo {
    margin-right: 3.7rem;
    margin-left: 0.75rem;
  }
}
</style>
