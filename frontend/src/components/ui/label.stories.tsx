import type { Meta, StoryObj } from '@storybook/react-vite';
import { expect } from 'storybook/test';
import { Label } from './label';
import { Input } from './input';

const meta = {
  component: Label,
  tags: ['ai-generated'],
} satisfies Meta<typeof Label>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = { args: { children: 'Email address' } };

export const AssociatedWithInput: Story = {
  render: () => (
    <div className="flex flex-col gap-1.5">
      <Label htmlFor="email-field">Email address</Label>
      <Input id="email-field" placeholder="you@example.com" />
    </div>
  ),
  play: async ({ canvas }) => {
    // Clicking the label focuses the associated input — proves `htmlFor` wiring.
    await canvas.getByText('Email address').click();
    await expect(canvas.getByPlaceholderText('you@example.com')).toHaveFocus();
  },
};
