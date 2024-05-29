package testfactory

import (
	crawlJobRepo "crawlquery/api/crawl/job/repository/mem"
	crawlLogRepo "crawlquery/api/crawl/log/repository/mem"
	crawlService "crawlquery/api/crawl/service"
	"crawlquery/api/domain"

	nodeRepo "crawlquery/api/node/repository/mem"
	nodeService "crawlquery/api/node/service"

	linkRepo "crawlquery/api/link/repository/mem"
	linkService "crawlquery/api/link/service"

	pageRepo "crawlquery/api/page/repository/mem"
	pageService "crawlquery/api/page/service"

	shardRepo "crawlquery/api/shard/repository/mem"
	shardService "crawlquery/api/shard/service"

	eventService "crawlquery/api/event/service"

	"crawlquery/pkg/testutil"
)

type ServiceFactory struct {
	EventService *eventService.Service
	ShardRepo    *shardRepo.Repository
	ShardService *shardService.Service
	PageRepo     *pageRepo.Repository
	PageService  *pageService.Service
	LinkRepo     *linkRepo.Repository
	LinkService  *linkService.Service
	NodeRepo     *nodeRepo.Repository
	NodeService  *nodeService.Service
	CrawlJobRepo *crawlJobRepo.Repository
	CrawlLogRepo *crawlLogRepo.Repository
	CrawlService *crawlService.Service
}

type ServiceFactoryOption func(*ServiceFactory)

func WithShard(shard *domain.Shard) ServiceFactoryOption {
	return func(factory *ServiceFactory) {
		factory.ShardRepo.Create(shard)
	}
}

func WithNode(node *domain.Node) ServiceFactoryOption {
	return func(factory *ServiceFactory) {
		factory.NodeRepo.Create(node)
	}
}

func WithCrawlJob(job *domain.CrawlJob) ServiceFactoryOption {
	return func(factory *ServiceFactory) {
		factory.CrawlJobRepo.Save(job)
	}
}

func NewServiceFactory(options ...ServiceFactoryOption) *ServiceFactory {

	eventService := eventService.NewService()

	shardRepo := shardRepo.NewRepository()
	shardService := shardService.NewService(
		shardService.WithRepo(shardRepo),
		shardService.WithLogger(testutil.NewTestLogger()),
	)

	nodeRepo := nodeRepo.NewRepository()
	nodeService := nodeService.NewService(
		nodeService.WithNodeRepo(nodeRepo),
		nodeService.WithLogger(testutil.NewTestLogger()),
		nodeService.WithShardService(shardService),
	)

	linkRepo := linkRepo.NewRepository()
	linkService := linkService.NewService(
		linkService.WithEventService(eventService),
		linkService.WithLinkRepo(linkRepo),
		linkService.WithLogger(testutil.NewTestLogger()),
	)

	crawlRepo := crawlJobRepo.NewRepository()
	crawlLogRepo := crawlLogRepo.NewRepository()
	crawlService := crawlService.NewService(
		crawlService.WithCrawlJobRepo(crawlRepo),
		crawlService.WithNodeService(nodeService),
		crawlService.WithLinkService(linkService),
		crawlService.WithCrawlLogRepo(
			crawlLogRepo,
		),
		crawlService.WithLogger(testutil.NewTestLogger()),
	)

	pageRepo := pageRepo.NewRepository()
	pageService := pageService.NewService(
		pageService.WithEventService(eventService),
		pageService.WithPageRepo(pageRepo),
		pageService.WithLogger(testutil.NewTestLogger()),
		pageService.WithShardService(shardService),
	)

	factory := &ServiceFactory{
		EventService: eventService,
		ShardRepo:    shardRepo,
		ShardService: shardService,
		PageRepo:     pageRepo,
		PageService:  pageService,
		LinkRepo:     linkRepo,
		LinkService:  linkService,
		NodeRepo:     nodeRepo,
		NodeService:  nodeService,
		CrawlJobRepo: crawlRepo,
		CrawlLogRepo: crawlLogRepo,
		CrawlService: crawlService,
	}
	for _, option := range options {
		option(factory)
	}
	return factory
}
