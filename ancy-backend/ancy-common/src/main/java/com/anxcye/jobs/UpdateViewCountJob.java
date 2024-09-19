package com.anxcye.jobs;

import com.anxcye.service.ArticleService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

@Component
public class UpdateViewCountJob {

    @Autowired
    private ArticleService articleService;

    @Scheduled(cron = "0 */10 * * * ?")
    public void updateViewCount() {
        articleService.syncFromRedisToDB();
    }
    
}
