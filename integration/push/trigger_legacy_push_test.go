package push

import (
	"path/filepath"

	"code.cloudfoundry.org/cli/integration/helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("triggering legacy push", func() {
	var (
		appName       string
		host          string
		defaultDomain string
		privateDomain string
	)

	BeforeEach(func() {
		appName = helpers.NewAppName()
		host = helpers.NewAppName()
		defaultDomain = defaultSharedDomain()

		privateDomain = helpers.DomainName()
		domain := helpers.NewDomain(organization, privateDomain)
		domain.Create()
	})

	Context("when there are global properties in the manifest", func() {
		It("triggering old push with deprecation warning", func() {
			helpers.WithHelloWorldApp(func(dir string) {
				helpers.WriteManifest(filepath.Join(dir, "manifest.yml"), map[string]interface{}{
					"host": host,
					"applications": []map[string]string{
						{
							"name": appName,
						},
					},
				})

				session := helpers.CustomCF(helpers.CFEnv{WorkingDirectory: dir}, PushCommandName, "--no-start")
				Eventually(session.Err).Should(Say("App manifest has attributes promoted to the top level\\. Found: host\\."))
				Eventually(session).Should(Say("Creating route %s\\.%s", host, defaultDomain))
			})
		})
	})

	Context("when there is an 'inherit' property in the manifest", func() {
		It("triggering old push with deprecation warning", func() {
			helpers.WithHelloWorldApp(func(dir string) {
				helpers.WriteManifest(filepath.Join(dir, "parent.yml"), map[string]interface{}{
					"host": host,
				})

				helpers.WriteManifest(filepath.Join(dir, "manifest.yml"), map[string]interface{}{
					"inherit": "./parent.yml",
					"applications": []map[string]string{
						{
							"name": appName,
						},
					},
				})

				session := helpers.CustomCF(helpers.CFEnv{WorkingDirectory: dir}, PushCommandName, "--no-start")
				Eventually(session.Err).Should(Say("App manifest uses inheritance, which is not supported by 'v2-push'"))
				Eventually(session).Should(Say("Creating route %s\\.%s", host, defaultDomain))
				Eventually(session).Should(Say("OK"))
				Eventually(session).Should(Exit(0))
			})
		})
	})

	Context("when there is a 'domain' property in the manifest", func() {
		It("triggering old push with deprecation warning", func() {
			helpers.WithHelloWorldApp(func(dir string) {
				helpers.WriteManifest(filepath.Join(dir, "manifest.yml"), map[string]interface{}{
					"applications": []map[string]string{
						{
							"name":   appName,
							"domain": defaultDomain,
						},
					},
				})

				session := helpers.CustomCF(helpers.CFEnv{WorkingDirectory: dir}, PushCommandName, "--no-start")
				Eventually(session.Err).Should(Say("App manifest declares routes using domain or domains attributes."))
				Eventually(session.Err).Should(Say("These attributes are not processed by 'v2-push' and may be deprecated in the future."))
				Eventually(session).Should(Say("(?i)Creating route %s\\.%s", appName, defaultDomain))
				Eventually(session).Should(Say("OK"))
				Eventually(session).Should(Exit(0))
			})
		})
	})

	Context("when there is a 'domains' property in the manifest", func() {
		It("triggering old push with deprecation warning", func() {
			helpers.WithHelloWorldApp(func(dir string) {
				helpers.WriteManifest(filepath.Join(dir, "manifest.yml"), map[string]interface{}{
					"applications": []map[string]interface{}{
						{
							"name":    appName,
							"domains": []string{defaultDomain},
						},
					},
				})

				session := helpers.CustomCF(helpers.CFEnv{WorkingDirectory: dir}, PushCommandName, "--no-start")
				Eventually(session.Err).Should(Say("App manifest declares routes using domain or domains attributes."))
				Eventually(session.Err).Should(Say("These attributes are not processed by 'v2-push' and may be deprecated in the future."))
				Eventually(session).Should(Say("(?i)Creating route %s\\.%s", appName, defaultDomain))
				Eventually(session).Should(Say("OK"))
				Eventually(session).Should(Exit(0))
			})
		})
	})

	Context("when there is a 'host' property in the manifest", func() {
		It("triggering old push with deprecation warning", func() {
			helpers.WithHelloWorldApp(func(dir string) {
				helpers.WriteManifest(filepath.Join(dir, "manifest.yml"), map[string]interface{}{
					"applications": []map[string]string{
						{
							"name": appName,
							"host": host,
						},
					},
				})

				session := helpers.CustomCF(helpers.CFEnv{WorkingDirectory: dir}, PushCommandName, "--no-start")
				Eventually(session.Err).Should(Say("App manifest declares routes using 'host', 'hosts', or 'no-hostname'"))
				Eventually(session.Err).Should(Say("These attributes are not processed by 'v2-push'"))
				Eventually(session).Should(Say("(?i)Creating route %s\\.%s", host, defaultDomain))
				Eventually(session).Should(Say("OK"))
				Eventually(session).Should(Exit(0))
			})
		})
	})

	Context("when there is a 'hosts' property in the manifest", func() {
		It("triggering old push with deprecation warning", func() {
			helpers.WithHelloWorldApp(func(dir string) {
				helpers.WriteManifest(filepath.Join(dir, "manifest.yml"), map[string]interface{}{
					"applications": []map[string]interface{}{
						{
							"name":  appName,
							"hosts": []string{host},
						},
					},
				})

				session := helpers.CustomCF(helpers.CFEnv{WorkingDirectory: dir}, PushCommandName, "--no-start")
				Eventually(session.Err).Should(Say("App manifest declares routes using 'host', 'hosts', or 'no-hostname'"))
				Eventually(session.Err).Should(Say("These attributes are not processed by 'v2-push'"))
				Eventually(session).Should(Say("(?i)Creating route %s\\.%s", host, defaultDomain))
				Eventually(session).Should(Say("OK"))
				Eventually(session).Should(Exit(0))
			})
		})
	})

	Context("when there is a 'no-hostname' property in the manifest", func() {

		It("triggering old push with deprecation warning", func() {
			helpers.WithHelloWorldApp(func(dir string) {
				helpers.WriteManifest(filepath.Join(dir, "manifest.yml"), map[string]interface{}{
					"applications": []map[string]interface{}{
						{
							"name":        appName,
							"no-hostname": true,
						},
					},
				})

				session := helpers.CustomCF(helpers.CFEnv{WorkingDirectory: dir}, PushCommandName, "--no-start", "-d", privateDomain)
				Eventually(session.Err).Should(Say("App manifest declares routes using 'host', 'hosts', or 'no-hostname'"))
				Eventually(session.Err).Should(Say("These attributes are not processed by 'v2-push'"))
				Eventually(session).Should(Say("(?i)Creating route %s", privateDomain))
				Eventually(session).Should(Say("OK"))
				Eventually(session).Should(Exit(0))
			})
		})
	})

	Context("when there is an 'inherit' property and a 'hostname' property in the manifest", func() {
		It("triggering old push with deprecation warning", func() {
			helpers.WithHelloWorldApp(func(dir string) {
				helpers.WriteManifest(filepath.Join(dir, "parent.yml"), map[string]interface{}{
					"domain": privateDomain,
				})

				helpers.WriteManifest(filepath.Join(dir, "manifest.yml"), map[string]interface{}{
					"inherit": "./parent.yml",
					"applications": []map[string]string{
						{
							"name": appName,
							"host": host,
						},
					},
				})

				session := helpers.CustomCF(helpers.CFEnv{WorkingDirectory: dir}, PushCommandName, "--no-start")
				Eventually(session.Err).Should(Say("App manifest uses inheritance, which is not supported by 'v2-push'"))
				Eventually(session).Should(Say("Creating route %s\\.%s", host, privateDomain))
				Eventually(session).Should(Say("OK"))
				Eventually(session).Should(Exit(0))
			})
		})
	})

	Context("when there is a 'domain' property and a 'hostname' property in the manifest", func() {
		It("triggering old push with deprecation warning", func() {
			helpers.WithHelloWorldApp(func(dir string) {
				helpers.WriteManifest(filepath.Join(dir, "manifest.yml"), map[string]interface{}{
					"applications": []map[string]string{
						{
							"name":   appName,
							"domain": privateDomain,
							"host":   host,
						},
					},
				})

				session := helpers.CustomCF(helpers.CFEnv{WorkingDirectory: dir}, PushCommandName, "--no-start")
				Eventually(session.Err).Should(Say("App manifest declares routes using domain or domains attributes."))
				Eventually(session.Err).Should(Say("These attributes are not processed by 'v2-push' and may be deprecated in the future."))
				Eventually(session).Should(Say("(?i)Creating route %s\\.%s", host, privateDomain))
				Eventually(session).Should(Say("OK"))
				Eventually(session).Should(Exit(0))
			})
		})
	})
})
