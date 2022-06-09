export default function Index() {
  return (
    <div className="rootWrapper">
      <h1>Remix - Golang Authentication Demo</h1>
      <div className="mt-20 flex gap-5">
        <a
          className="btnLink-base text-white font-bold bg-blue-500"
          href="/register"
        >
          Register
        </a>
        <a className="btnLink-base bg-gray-200" href="login">
          Login
        </a>
      </div>
    </div>
  );
}
