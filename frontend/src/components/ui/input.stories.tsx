import type { Meta, StoryObj } from '@storybook/react-vite';
import { expect } from 'storybook/test';
import { Input } from './input';

const meta = {
  component: Input,
  tags: ['ai-generated'],
} satisfies Meta<typeof Input>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = { args: { placeholder: 'you@example.com', type: 'email' } };
export const Disabled: Story = { args: { placeholder: 'Disabled', disabled: true } };

export const Typing: Story = {
  args: { placeholder: 'you@example.com', type: 'email' },
  play: async ({ canvas, userEvent }) => {
    const input = canvas.getByPlaceholderText('you@example.com');
    await userEvent.type(input, 'a@b.com');
    await expect(input).toHaveValue('a@b.com');
  },
};
