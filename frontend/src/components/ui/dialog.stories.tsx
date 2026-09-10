import type { Meta, StoryObj } from '@storybook/react-vite';
import { expect, within, waitFor } from 'storybook/test';
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from './dialog';
import { Button } from './button';

const meta = {
  component: Dialog,
  tags: ['ai-generated'],
} satisfies Meta<typeof Dialog>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <Dialog>
      <DialogTrigger asChild>
        <Button>Cancel transaction</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Cancel this transaction?</DialogTitle>
          <DialogDescription>This action cannot be undone.</DialogDescription>
        </DialogHeader>
        <DialogFooter showCloseButton>
          <Button variant="destructive">Confirm</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  ),
  // DialogContent renders into a portal, so it lands outside the story's own
  // canvas root — query the trigger via `canvas` but the portal content via
  // the owner document's body.
  play: async ({ canvas, userEvent, canvasElement }) => {
    await userEvent.click(canvas.getByRole('button', { name: /cancel transaction/i }));
    const body = within(canvasElement.ownerDocument.body);
    await body.findByText('Cancel this transaction?');
    // The dialog fades in (data-open:animate-in), so wait for the animation
    // to finish rather than asserting visibility on first DOM appearance.
    await waitFor(() => expect(body.getByText('Cancel this transaction?')).toBeVisible());
  },
};
