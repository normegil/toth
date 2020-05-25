<template>
  <div class="main-menu-item" :class="{'main-menu-item--bottom': isBottom, 'main-menu-item--top': !isBottom}" @click="click">
      <slot class="main-menu-item__icon"></slot>
      <span class="main-menu-item__text"> {{ title }}</span>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";

@Component
export default class MainMenuItem extends Vue {
  @Prop({ required: true })
  title!: string;

  @Prop({ default: false, required: false })
  isBottom!: boolean

  click() {
      this.$emit("click");
  }
}
</script>

<style lang="scss">
.main-menu-item {
  display: block;
  font-size: 1.9rem;
  padding: 0.75rem 0;

  padding-left: 2.2rem;
  color: $color-primary-light;
  text-decoration: none;
  cursor: pointer;
  transition: all 0.3s;

  &:hover {
    background-color: $color-complementary-lighter;
    color: $color-complementary;
  }

  &__text {
    display: inline-block;
    margin-left: 0.75rem;
  }
    &--top {
        border-bottom: 1px solid $color-primary-light;
        &:hover {
            border-bottom: 1px solid $color-complementary;
        }
    }
    &--bottom {
        border-top: 1px solid $color-primary-light;
        &:hover {
            border-top: 1px solid $color-complementary;
        }
    }
}
</style>
