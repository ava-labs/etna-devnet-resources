import NextPrev from "./ui/NextPrev";

export default function LaunchRpcNode() {
    return (
        <div>
            <h1 className="text-2xl font-medium mb-6">Launch RPC Node</h1>

            {/* ... existing/future content ... */}

            <NextPrev
                nextDisabled={false}
                currentStepName="launch-rpc-node"
            />
        </div>
    );
}
