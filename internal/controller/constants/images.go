package constants

var (
	TrillianLogSignerImage = "quay.io/securesign/trillian-logsigner@sha256:37028258a88bba4dfaadb59fc88b6efe9c119a808e212ad5214d65072abb29d0"
	TrillianServerImage    = "quay.io/securesign/trillian-logserver@sha256:994a860e569f2200211b01f9919de11d14b86c669230184c4997f3d875c79208"
	TrillianDbImage        = "quay.io/securesign/trillian-database@sha256:909f584804245f8a9e05ecc4d6874c26d56c0d742ba793c1a4357a14f5e67eb0"

	// TODO: remove and check the DB pod status
	OseToolsImage = "registry.redhat.io/openshift4/ose-tools-rhel8@sha256:486b4d2dd0d10c5ef0212714c94334e04fe8a3d36cf619881986201a50f123c7"

	FulcioServerImage = "quay.io/securesign/fulcio-server@sha256:67495de82e2fcd2ab4ad0e53442884c392da1aa3f5dd56d9488a1ed5df97f513"

	RekorRedisImage    = "quay.io/securesign/trillian-redis@sha256:01736bdd96acbc646334a1109409862210e5273394c35fb244f21a143af9f83e"
	RekorServerImage   = "quay.io/securesign/rekor-server@sha256:133ee0153e12e6562cfea1a74914ebdd7ee76ae131ec7ca0c3e674c2848150ae"
	RekorSearchUiImage = "quay.io/securesign/rekor-search-ui@sha256:8c478fc6122377c6c9df0fddf0ae42b6f6b1648e3c6cf96a0558f366e7921b2b"
	BackfillRedisImage = "quay.io/securesign/rekor-backfill-redis@sha256:88869eb582cbb94baa50c212689c50ed405cc94669c2c03f781b12ad867827ce"

	TufImage = "quay.io/securesign/scaffold-tuf-server@sha256:34f5cdc53a908ae2819d85ab18e35b69dc4efc135d747dd1d2e216a99a2dcd1b"

	CTLogImage = "quay.io/securesign/certificate-transparency-go@sha256:a0c7d71fc8f4cb7530169a6b54dc3a67215c4058a45f84b87bb04fc62e6e8141"

	ClientServerImage       = "registry.access.redhat.com/ubi9/httpd-24@sha256:7874b82335a80269dcf99e5983c2330876f5fe8bdc33dc6aa4374958a2ffaaee"
	ClientServerImage_cg    = "quay.io/securesign/cli-client-server-cg@sha256:987c630213065a6339b2b2582138f7b921473b86dfe82e91a002f08386a899ed"
	ClientServerImage_re    = "quay.io/securesign/client-server-re@sha256:dc4667af49ce6cc70d70bf83cab9d7a14b424d8ae1aae7e4863ff5c4ac769a96"
	ClientServerImage_f     = "quay.io/securesign/client-server-f@sha256:65fb59c8f631215d9752fc4f41571eb2750ecaaa8555083f58baa6982e97d192"
	SegmentBackupImage      = "quay.io/securesign/segment-backup-job@sha256:3fcf8f14a0cfdd36f9ec263f83ba1597f892e6fa923d3d61bacbc467af643c9d"
	TimestampAuthorityImage = "quay.io/securesign/timestamp-authority@sha256:3fba2f8cd09548d2bd2dfff938529952999cb28ff5b7ea42c1c5e722b8eb827f"
)
