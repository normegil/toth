<template>
  <div class="logo__container">
    <LogoNoText class="logo__image" />
    <div class="logo__text">Toth</div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import LogoNoText from "../../assets/images/logo-notext.svg";

@Component({
  components: { LogoNoText }
})
export default class Logo extends Vue {}
</script>

<style lang="scss">
.logo {
  &__container {
    padding: 0.5rem;
    cursor: pointer;
  }
  &__image {
    display: inline-block;
    position: relative;
    top: 0.65rem;
    height: 3rem;
    fill: $color-grey-6;
    transition: all 0.2s;
  }
  &__text {
    display: inline-block;
    font-family: "Anke Calligraphic FG", sans-serif;
    font-size: 2.5rem;
    color: $color-grey-6;
    transition: all 0.2s;
  }

  &__container:hover {
    > .logo__image {
      fill: $color-complementary-light;
    }
    > .logo__text {
      color: $color-complementary-light;
    }
  }

  &__container:active {
    > .logo__image {
      fill: $color-complementary;
    }
    > .logo__text {
      color: $color-complementary;
    }
  }
}
</style>
