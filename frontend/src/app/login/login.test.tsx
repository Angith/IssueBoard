import { render, screen, fireEvent } from '@testing-library/react';
import LoginPage from './page';
import { supabase } from '@/lib/supabase';

jest.mock('@/lib/supabase', () => ({
  supabase: {
    auth: {
      signInWithOtp: jest.fn(),
    },
  },
}));

describe('LoginPage', () => {
  it('calls signInWithOtp on form submission', async () => {
    (supabase.auth.signInWithOtp as jest.Mock).mockResolvedValue({ error: null });
    
    render(<LoginPage />);
    
    const emailInput = screen.getByPlaceholderText('Your email');
    const submitButton = screen.getByText('Send Magic Link');
    
    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
    fireEvent.click(submitButton);
    
    expect(supabase.auth.signInWithOtp).toHaveBeenCalledWith({
      email: 'test@example.com',
      options: {
        emailRedirectTo: expect.stringContaining('/inventory'),
      },
    });
  });
});
