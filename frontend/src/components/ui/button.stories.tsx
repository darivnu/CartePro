import type { Meta, StoryObj } from '@storybook/react-vite';
import { expect } from 'storybook/test';
import { Button } from './button';

const meta = {
  component: Button,
  tags: ['ai-generated'],
} satisfies Meta<typeof Button>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = { args: { children: 'Continue' } };
export const Outline: Story = { args: { children: 'Continue', variant: 'outline' } };
export const Secondary: Story = { args: { children: 'Continue', variant: 'secondary' } };
export const Destructive: Story = { args: { children: 'Delete', variant: 'destructive' } };

export const Disabled: Story = {
  args: { children: 'Continue', disabled: true },
  play: async ({ canvas }) => {
    await expect(canvas.getByRole('button', { name: /continue/i })).toBeDisabled();
  },
};

// Button's default variant uses bg-primary — fails if Tailwind / the app's
// index.css theme did not load into the Storybook preview.
export const CssCheck: Story = {
  args: { children: 'Continue' },
  play: async ({ canvas }) => {
    const button = canvas.getByRole('button', { name: /continue/i });
    await expect(getComputedStyle(button).backgroundColor).not.toBe('rgba(0, 0, 0, 0)');
  },
};
