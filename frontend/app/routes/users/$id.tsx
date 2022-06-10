import type { LoaderFunction } from '@remix-run/node';
import { useLoaderData } from '@remix-run/react';

export const loader: LoaderFunction = async ({ params }) => {
  const;
};

export default function UserPage() {
  const { name } = useLoaderData();

  return (
    <div className="rootWrapper">
      <h1>Welcome home, {name}</h1>
      <button className="btnLink-base">Logout</button>
      <button className="btnLink-base">Delete account</button>
    </div>
  );
}
