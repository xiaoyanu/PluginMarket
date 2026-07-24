<script setup lang="ts">

type ButtonType = 'blue' | 'red' | 'orange' | 'green' | 'gray' | 'black' | 'white'
type ButtonSize = 'large' | 'small'

interface Props {
  text?: string
  color?: ButtonType
  size?: ButtonSize
  disabled?: boolean
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  text: '',
  color: 'blue',
  size: undefined, // 默认不传即为标准/中等尺寸
  disabled: false,
  loading: false
})

const isDisabled = computed(() => props.disabled || props.loading)
</script>

<template>
  <button
      type="button"
      class="pm-button"
      :class="[
        `btn-${color}`,
        size ? `btn-${size}` : '',
        { 'is-loading': loading }
      ]"
      :disabled="isDisabled"
  >
    <span v-if="loading" class="loading-icon">
      <svg viewBox="25 25 50 50" class="circular">
        <circle cx="50" cy="50" r="20" fill="none" class="path"></circle>
      </svg>
    </span>

    <slot v-if="!loading" name="icon"/>

    <span v-if="text || $slots.default" class="text-content">
      {{ text }}
    </span>
  </button>
</template>

<style scoped lang="scss">
.pm-button {
  display: inline-flex;
  justify-content: center;
  align-items: center;
  gap: 6px;

  padding: 8px 15px;
  font-size: 14px;
  border-radius: 7px;

  font-weight: 500;
  line-height: 1;
  color: #fff;

  white-space: nowrap;
  user-select: none;
  cursor: pointer;
  appearance: none;
  border: 1px solid transparent;
  transition: background-color 0.3s cubic-bezier(0.4, 0, 0.2, 1),
  box-shadow 0.3s cubic-bezier(0.4, 0, 0.2, 1),
  transform 0.1s ease-out;

  &:active:not(:disabled) {
    transform: translateY(1px);
  }

  &:disabled {
    cursor: not-allowed;
    opacity: 0.6;
  }
}

.btn-small {
  padding: 5px 10px;
  font-size: 12px;
  border-radius: 5px;
}

.btn-large {
  padding: 12px 20px;
  font-size: 16px;
  border-radius: 9px;
}

@mixin button-variant($bg-color, $shadow-color, $hover-color, $active-color) {
  background-color: $bg-color;
  box-shadow: 0 2px 8px rgba($shadow-color, 0.35);

  &:hover:not(:disabled) {
    background-color: $hover-color;
    box-shadow: 0 4px 12px rgba($shadow-color, 0.45);
  }

  &:active:not(:disabled) {
    background-color: $active-color;
    box-shadow: 0 2px 4px rgba($shadow-color, 0.2);
  }
}

/* 颜色变体 */
.btn-blue {
  @include button-variant(#4E89FF, #4E89FF, #3B76EF, #2D65D9);
}

.btn-red {
  @include button-variant(#FF6B6B, #FF6B6B, #FA5252, #F03E3E);
}

.btn-orange {
  @include button-variant(#FFA94D, #FFA94D, #FF922B, #FD7E14);
}

.btn-green {
  @include button-variant(#51CF66, #51CF66, #40C057, #37B24D);
}

.btn-gray {
  @include button-variant(#ADB5BD, #ADB5BD, #868E96, #495057);
}

.btn-black {
  @include button-variant(#343A40, #212529, #212529, #000000);
}

.btn-white {
  @include button-variant(#FFFFFF, #E2E8F0, #F8FAFC, #F1F5F9);
  border: 1px solid #DCDFE6;
  color: #1A1A1A;
}

/* Loading 动画 */
.loading-icon {
  width: 1em;
  height: 1em;
  display: flex;
  align-items: center;
  justify-content: center;

  .circular {
    display: block;
    width: 100%;
    height: 100%;
    animation: rotate 2s linear infinite;
  }

  .path {
    stroke: currentColor;
    stroke-width: 5;
    stroke-linecap: round;
    animation: dash 1.5s ease-in-out infinite;
  }
}

@keyframes rotate {
  100% {
    transform: rotate(360deg);
  }
}

@keyframes dash {
  0% {
    stroke-dasharray: 1, 150;
    stroke-dashoffset: 0;
  }
  50% {
    stroke-dasharray: 90, 150;
    stroke-dashoffset: -35;
  }
  100% {
    stroke-dasharray: 90, 150;
    stroke-dashoffset: -124;
  }
}
</style>