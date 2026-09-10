import type { Meta, StoryObj } from '@storybook/react-vite';
import { expect } from 'storybook/test';
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from './card';
import { Button } from './button';

const meta = {
  component: Card,
  tags: ['ai-generated'],
} satisfies Meta<typeof Card>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: (args) => (
    <Card {...args} className="w-80">
      <CardHeader>
        <CardTitle>Partner offer</CardTitle>
        <CardDescription>10% off your next visit</CardDescription>
      </CardHeader>
      <CardContent>Show your card at checkout to redeem.</CardContent>
      <CardFooter>
        <Button size="sm">Redeem</Button>
      </CardFooter>
    </Card>
  ),
  play: async ({ canvas }) => {
    await expect(canvas.getByText('Partner offer')).toBeVisible();
  },
};

export const Small: Story = {
  render: (args) => (
    <Card {...args} size="sm" className="w-80">
      <CardHeader>
        <CardTitle>Compact card</CardTitle>
      </CardHeader>
      <CardContent>Uses the reduced --card-spacing scale.</CardContent>
    </Card>
  ),
};
