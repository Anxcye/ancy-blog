package com.anxcye.jobs;

import com.anxcye.constants.RedisConstant;
import com.anxcye.service.ArticleService;
import com.anxcye.utils.RedisCache;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

import java.util.Map;

@Component
public class UpdateViewCountJob {

    @Autowired
    private RedisCache redisCache;

    @Autowired
    private ArticleService articleService;

    @Scheduled(cron = "0/5 * * * * ?")
    public void updateViewCount() {
        Map<String, Integer> viewCountMap = redisCache.getCacheMap(RedisConstant.ARTICLE_VIEW_COUNT);
        articleService.updateViewCount(viewCountMap);
    }
    
}
