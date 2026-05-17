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
    <div className="flex flex-col items-center justify-center min-h-screen bg-[#09090b] px-4 selection:bg-zinc-800">
      <div className="w-full max-w-sm p-8 rounded-2xl bg-zinc-900/20 backdrop-blur-md border border-zinc-800/50 shadow-2xl">
        <div className="mb-8 text-center">
          <h1 className="text-2xl font-medium tracking-tight text-zinc-100">Log in to IssueBoard</h1>
          <p className="mt-2 text-sm text-zinc-400">Enter your email to receive a magic link</p>
        </div>
        <form onSubmit={handleLogin} className="flex flex-col gap-5">
          <div className="flex flex-col gap-1.5">
            <label htmlFor="email" className="text-xs font-medium text-zinc-400 ml-1">
              Email address
            </label>
            <input
              id="email"
              type="email"
              placeholder="name@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full p-2.5 bg-zinc-900/50 border border-zinc-800 rounded-lg text-sm text-zinc-100 placeholder:text-zinc-600 focus:outline-none focus:ring-1 focus:ring-zinc-500 focus:border-zinc-500 transition-all duration-200"
              required
            />
          </div>
          <button
            type="submit"
            disabled={loading}
            className="w-full p-2.5 bg-zinc-100 text-zinc-900 font-medium text-sm rounded-lg hover:bg-white focus:outline-none focus:ring-2 focus:ring-zinc-500 focus:ring-offset-2 focus:ring-offset-[#09090b] transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? 'Sending...' : 'Continue with Email'}
          </button>
        </form>
        {message && (
          <div className="mt-6 p-3 rounded-md bg-zinc-900/50 border border-zinc-800/50">
            <p className="text-sm text-center text-zinc-300">{message}</p>
          </div>
        )}
      </div>
    </div>
  );
}
