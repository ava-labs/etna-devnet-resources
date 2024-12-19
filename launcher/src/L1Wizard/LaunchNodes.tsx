
import NextPrev from './ui/NextPrev';

export default function LaunchNodes() {

    return <>
        <h1 className="text-2xl font-medium mb-6">Launch L1 Nodes</h1>

        <NextPrev nextDisabled={false} currentStepName="generate-keys" />
    </>
}
