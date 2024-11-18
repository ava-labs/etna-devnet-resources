export default function Card({ children, title }: { children: React.ReactNode, title: string }) {
    return <div className="bg-white rounded-lg shadow-md p-8 m-8">
        <h2 className="text-2xl font-bold pb-4">{title}</h2>
        {children}
    </div>
}

