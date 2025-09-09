'use client';
import { useRouter,useSearchParams } from 'next/navigation';
import { Suspense, useEffect, useState } from 'react';
import { useAuthStore } from '@/app/store/auth';
import Link from 'next/link'
// import OtpTimer from '@/app/ui/components/login/otptimer';


function LoginFormContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const {signinUser} = useAuthStore();

  const [username,setUsername] = useState('');
    const [password,setPassword] = useState('');
    const [pin,setPin] = useState('');
    const [error,setError] = useState('');
    const [loading, setLoading] = useState(false);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    try {
        const response = await signinUser(username,password,pin);
        
        router.replace('/');
    } catch (err) {
        setError('Login failed. Please check your credentials.');
        console.error(err);
    } finally {
        setLoading(false);
    }
  }



  return (
    <form action="#" method="POST" className="space-y-6">
    <div> ̰
      <label htmlFor="username" className="block text-sm font-medium leading-6 text-gray-900">
        Username
      </label>
      <div className="mt-2">
        <input
          id="username"
          name="username"
          type="text"
          required
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          autoComplete="username"
          className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary sm:text-sm sm:leading-6"
        />
      </div>
    </div>
    <div>
      <label htmlFor="pin" className="block text-sm font-medium leading-6 text-gray-900">
        PIN
      </label>
      <div className="mt-2">
        <input
          id="pin"
          name="pin"
          type="number"
          required
          value={pin}
          onChange={(e) => setPin(e.target.value)}
          autoComplete="otp"
          className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary sm:text-sm sm:leading-6"
        />
      </div>
    </div>
    <div>
      <label htmlFor="password" className="block text-sm font-medium leading-6 text-gray-900">
        Password
      </label>
      <div className="mt-2">
        <input
          id="password"
          name="password"
          type="text"
          required
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          autoComplete="password"
          className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary sm:text-sm sm:leading-6"
        />
      </div>
    </div>

    <div className="flex items-center justify-between">
      <div className="flex items-center">
        <input
          id="remember-me"
          name="remember-me"
          type="checkbox"
          className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
        />
        <label htmlFor="remember-me" className="ml-3 block text-sm leading-6 text-gray-900">
          Remember me
        </label>
      </div>
    </div>

    <div>
      <button
        type="button"
        onClick={(e) => onSubmit(e)}
        className="flex w-full justify-center rounded-md bg-gray-800 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
      >
        Sign in
      </button>
    </div>
  </form>

  )
}

function LoginForm() {
  return (
    <Suspense>
      <LoginFormContent />
    </Suspense>
  )
}

export default function Page() {
  const router = useRouter();
  const [checkedToken, setCheckedToken] = useState(false);
  

  const goHome = () => {
    router.replace('/');
  }
  function checkAccessToken() {
    const accessToken = localStorage.getItem('access_token');
    if (accessToken && accessToken.length > 0) {
      goHome();
    }
    setCheckedToken(true);
  }

  useEffect(() => {
    if (!checkedToken) {
      checkAccessToken();
    }   

  },[checkedToken]);


  return (
      <>
        <div className="flex min-h-full flex-1 flex-col justify-center py-12 sm:px-6 lg:px-8">
          <div className="sm:mx-auto sm:w-full sm:max-w-md">
            <h2 className="mt-6 text-center text-2xl font-bold leading-9 tracking-tight text-gray-100">
              Sign in to your account
            </h2>
          </div>
  
          <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-[480px]">
            <div className="bg-white px-6 py-12 shadow sm:rounded-lg sm:px-12">
              <Suspense fallback={<div>Loading...</div>}>
                <LoginForm />
              </Suspense>

              <div>
                <div className="relative mt-10">
                  <div aria-hidden="true" className="absolute inset-0 flex items-center">
                    <div className="w-full border-t border-gray-200" />
                  </div>
                  <div className="relative flex justify-center text-sm font-medium leading-6">
                    <span className="bg-white px-6 text-gray-900"><Link href="/">Back Home</Link></span>
                  </div>
                </div>
              </div>
            </div>
  
          </div>
        </div>
      </>
    )
  }
  