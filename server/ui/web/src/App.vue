<template>
  <div class="app">
    <Navigation class="app__navigation" />
    <div class="app__page">
      <TitleBar class="app__title-bar" />
      <div class="app__content">Content</div>
    </div>
    <AuthenticationScreen />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Navigation from "./layout/Navigation.vue";
import TitleBar from "./layout/title/TitleBar.vue";
import AuthenticationScreen from "./layout/login/AuthenticationScreen.vue";

@Component({
  components: { TitleBar, Navigation, AuthenticationScreen }
})
export default class App extends Vue {
  created(): void {
    this.$store.dispatch("auth/loadAuthentication");
  }
}
</script>

<style lang="scss">
.app {
  display: flex;
  height: 100vh;
  font-size: 1.6rem;
  overflow: hidden;

  &__navigation {
    min-width: 19rem;
  }

  &__page {
    width: 100%;
  }

  &__title-bar {
    min-height: 6.77rem;
  }

  &__content {
    background-color: $color-grey-7;
    height: 100%;
  }
}
</style>
