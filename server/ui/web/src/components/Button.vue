<template>
  <a class="button" @click="click">
    <i v-if="icon !== '' && icon != null" :class="icon" />
    <span v-if="title !== '' && title != null">{{ title }}</span>
  </a>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
@Component
export default class Button extends Vue {
  @Prop({ default: "", required: false })
  icon!: string;

  @Prop({ default: "", required: false })
  title!: string;

  click(): void {
    this.$emit("click");
  }
}
</script>

<style lang="scss">
.button {
  display: inline-block;
  padding: 0.5rem;
  border: 1px solid $color-grey-7;
  background-color: $color-grey-8;
  text-decoration: none;
  color: $color-grey-1;
  transition: all 0.3s;
  white-space: nowrap;
  cursor: pointer;

  &:active,
  &:hover {
    background-color: $color-complementary;
  }
}
</style>
