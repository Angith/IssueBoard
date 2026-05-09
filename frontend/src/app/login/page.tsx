'use client';

import { useState } from 'react';
import { supabase } from '@/lib/supabase';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    console.debug('[Login Debug] Login attempt started for email:', email);
    setLoading(true);
    setMessage('');

    console.debug('[Login Debug] Calling supabase.auth.signInWithOtp for:', email);
    const { error } = await supabase.auth.signInWithOtp({
      email,
      options: {
        emailRedirectTo: `${window.location.origin}/inventory`,
      },
    });

    if (error) {
      console.error('[Login Debug] Error during sign in:', error.message, error);
      setMessage(`Error: ${error.message}`);
    } else {
      console.info('[Login Debug] Magic link sent successfully. Awaiting user action.');
      setMessage('Check your email for the magic link!');
    }
    setLoading(false);
    console.debug('[Login Debug] Login attempt finished.');
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen py-2">
      <h1 className="text-4xl font-bold mb-8">Login to IssueBoard</h1>
      <form onSubmit={handleLogin} className="flex flex-col gap-4 w-full max-w-md">
        <input
          type="email"
          placeholder="Your email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="p-2 border rounded"
          required
        />
        <button
          type="submit"
          disabled={loading}
          className="p-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
        >
          {loading ? 'Sending...' : 'Send Magic Link'}
        </button>
      </form>
      {message && <p className="mt-4 text-center">{message}</p>}
    </div>
  );
}
