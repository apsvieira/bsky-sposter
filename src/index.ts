import "dotenv/config";
import Parser from "rss-parser";
import { AtpAgent, RichText } from "@atproto/api";

// Configuration
const RSS_FEED_URL = process.env.RSS_FEED_URL!;
const BLUESKY_HANDLE = process.env.BLUESKY_HANDLE!;
const BLUESKY_APP_PASSWORD = process.env.BLUESKY_APP_PASSWORD!;
const BLUESKY_SERVICE = process.env.BLUESKY_SERVICE || "https://bsky.social";

interface RSSItem {
  title: string;
  link: string;
  contentSnippet: string;
  guid: string;
  refDate: Date;
}

function ParseRSSItem(item: any): RSSItem {
  const rawDate = item.isoDate || item.pubDate || item.date;
  return {
    title: item.title,
    link: item.link,
    contentSnippet: item.contentSnippet,
    guid: item.guid,
    refDate: new Date(rawDate),
  };
}

function formatDate(date: Date): string {
  return date.toISOString().split("T")[0];
}

// Initialize BlueSky agent
const agent = new AtpAgent({ service: BLUESKY_SERVICE });

// Track processed posts to avoid duplicates
const processedPosts = new Set<string>();

// Function to fetch and parse RSS feed
async function fetchRSSFeed() {
  const parser = new Parser();
  const feed = await parser.parseURL(RSS_FEED_URL);
  const items = feed.items.map(ParseRSSItem);
  // Post in chronological order (older first)
  const sortedItems = items.sort(
    (a, b) => a.refDate.getTime() - b.refDate.getTime()
  );

  for (const item of sortedItems) {
    if (!item.guid || !item.title || !item.link) continue; // Skip if no GUID

    if (!processedPosts.has(item.guid)) {
      console.log(`Processing: ${item.title}`);
      const snippet = getFirstPhrase(item.contentSnippet || "");

      await postToBlueSky(
        item.title,
        formatDate(item.refDate),
        item.link,
        snippet
      );
      processedPosts.add(item.guid); // Mark as processed
    }
  }
}

function getFirstPhrase(content: string, maxLen: number = 200): string {
  if (!content) return "";

  const firstPhrase = content.split(".")[0];
  if (firstPhrase.length > maxLen) {
    return firstPhrase.slice(0, maxLen - 3) + "...";
  }
  return firstPhrase + "...";
}

// Function to post to BlueSky
async function postToBlueSky(
  title: string,
  date: string,
  link: string,
  content: string
) {
  await agent.login({
    identifier: BLUESKY_HANDLE,
    password: BLUESKY_APP_PASSWORD,
  });

  const richText = new RichText({
    text: `${title} (${date})\n\n${content}\n\n${link}`,
  });
  await richText.detectFacets(agent);

  await agent.app.bsky.feed.post.create(
    { repo: agent.session?.did }, // Use the DID from the session
    {
      text: richText.text,  
      facets: richText.facets,
      createdAt: new Date().toISOString(),
    }
  );

  console.log(`Posted: ${title}`);
}

setInterval(fetchRSSFeed, 10 * 60 * 1000); // Fetch every 10 minutes
fetchRSSFeed();
