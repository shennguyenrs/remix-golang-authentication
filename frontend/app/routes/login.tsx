import type { ActionFunction, ErrorBoundaryComponent } from '@remix-run/node';
import { redirect } from '@remix-run/node';
import { Form } from '@remix-run/react';
import axios from 'axios';

export const action: ActionFunction = async ({ request }) => {
  const form = await request.formData();
  const email = form.get('email');
  const password = form.get('password');

  try {
    const res = await axios.post(process.env.BACKEND_AUTH + '/login', {
      email,
      password,
    });

    if (res) {
      const { token, userId } = res.data;
			console.log(token)
      request.headers.set('Authorization', token);
      return redirect('/users/' + userId);
    }
  } catch (err) {
    if (axios.isAxiosError(err) && err.response) {
      throw new Error(err.response.data as string);
    }

    throw new Error('An error occurred while logging in');
  }
};

export const ErrorBoundary: ErrorBoundaryComponent = ({
  error,
}: {
  error: Error;
}) => {
  return (
    <div className="rootWrapper">
      <p>{error.message}</p>
      <a className="btnLink-base text-white bg-blue-500 mt-5" href="/">
        Return home
      </a>
    </div>
  );
};

export default function Login() {
  return (
    <div className="rootWrapper">
      <h1>Login page</h1>
      <Form className="formCon" method="post">
        <label>Email</label>
        <input type="email" name="email" />
        <label>Password</label>
        <input type="password" name="password" />
        <button
          className="btnLink-base font-bold text-white bg-blue-500"
          type="submit"
        >
          Login
        </button>
      </Form>
      <div>
        <a className="btnLink-base text-white bg-red-300" href="/">
          Go back to homepage
        </a>
      </div>
    </div>
  );
}
