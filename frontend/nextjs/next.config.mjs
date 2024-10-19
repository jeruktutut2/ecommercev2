/** @type {import('next').NextConfig} */
const nextConfig = {
    reactStrictMode: true,
    async rewrites() {
        return [
            {
                source: '/api/v1/users/:path*',
                destination: `http://localhost:10001/api/v1/users/:path*`,
            },
        ]
    },
};

export default nextConfig;
