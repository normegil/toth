<template>
  <div
    class="field__container"
    :class="{ 'field__container--disabled': disabled }"
  >
    <label
      :for="this._uid"
      class="field__label"
      :class="{ 'field__label--disabled': disabled }"
      v-if="labelEnabled"
      >{{ label }}</label
    >
    <input
      :id="this._uid"
      class="field__input"
      :class="{
        'left-rounded': !labelEnabled && !password,
        'right-rounded': !buttonEnabled && !password
      }"
      :value="value"
      :placeholder="placeholder"
      @input="$emit('input', $event.target.value)"
      @keyup.enter="$emit('keyup-enter')"
      :disabled="disabled"
      :type="fieldType"
    />
    <input
      v-if="password && enableShowPassword"
      type="checkbox"
      v-model="showPassword"
      class="field__password-checkbox"
      :id="this._uid + 'password-checkbox'"
    /><label
      v-if="password && enableShowPassword"
      :for="this._uid + 'password-checkbox'"
      class="field__password-checkbox-label"
      :class="{
        'right-rounded': !buttonEnabled,
      }"
    ><span v-if="enableShowPassword && showPassword"><EyeSlash class="field__password-checkbox-icon"/></span><span v-if="!showPassword" ><Eye class="field__password-checkbox-icon field__password-checkbox-icon--eye " /></span></label>
    <a
      v-if="buttonEnabled"
      class="field__button"
      @click.stop="$emit('button-click')"
    >
      {{ button }}
    </a>
  </div>
</template>

<script lang="ts">
import Component from "vue-class-component";
import Vue from "vue";
import { Prop } from "vue-property-decorator";
import Eye from "../assets/images/icons/eye-regular.svg"
import EyeSlash from "../assets/images/icons/eye-slash-regular.svg"

@Component({
  components: {
    Eye,
    EyeSlash
  }
})
export default class InputField extends Vue {
  @Prop({ default: "", required: false })
  label!: string;

  @Prop({ required: true })
  value!: string;

  @Prop({ default: "", required: false })
  placeholder!: string;

  @Prop({ default: "", required: false })
  button!: string;

  @Prop({ default: false, required: false })
  disabled!: boolean;

  @Prop({ default: false, required: false })
  password!: boolean;

  @Prop({ default: true, required: false })
  enableShowPassword!: boolean;

  showPassword = false;

  get fieldType(): string {
    if (this.password && !this.showPassword) {
      return "password";
    }
    return "text";
  }

  get labelEnabled(): boolean {
    return this.label !== "";
  }

  get buttonEnabled(): boolean {
    return this.button !== "";
  }
}
</script>

<style lang="scss">
.field {
  &__container {
    display: flex;
    flex-direction: row;
    border: 1px solid $color-grey-7;
    font-size: 1.6rem;
    transition: all 0.3s;

    &:hover:not(&--disabled),
    &:focus:not(&--disabled) {
      border: 1px solid $color-complementary;

      & > .field__label {
        background: $color-complementary;
      }
    }

    &--disabled {
      cursor: not-allowed;
    }
  }
  &__label {
    text-align: center;
    padding: 0.5rem 2rem;
    border-right: 1px solid $color-grey-8;
    background-color: $color-grey-8;
    transition: all 0.3s;

    &--disabled {
      cursor: not-allowed;
    }
  }
  &__input {
    flex-grow: 2;
    border: none;

    color: $color-grey-1;
    transition: all 0.3s;
    padding: 0 1rem;

    &:hover,
    &:focus {
      outline: none;
    }

    &--disabled {
      cursor: not-allowed;
    }
  }

  &__button {
    text-align: center;
    padding: 0.5rem 2rem;
    border-left: 1px solid $color-grey-7;
    transition: all 0.3s;

    &:hover,
    &:active {
      background: $color-grey-7;
    }
  }

  &__password-checkbox {
    display: none;
  }

  &__password-checkbox-label {
    background-color: $color-grey-9;
    padding: 0.6rem 1rem;
    cursor: pointer;

    &:hover,
    &:active {
      background-color: $color-complementary-lighter;
    }
  }

  &__password-checkbox-icon {
    width: 1.5rem;

    &--eye {
      position: relative;
      top: 1px;
    }
  }
}
</style>
