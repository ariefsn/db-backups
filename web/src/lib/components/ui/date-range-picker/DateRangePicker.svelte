<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Popover from '$lib/components/ui/popover';
	import { cn } from '$lib/utils';
	import { CalendarDate, DateFormatter, getLocalTimeZone } from '@internationalized/date';
	import { Calendar as CalendarIcon } from '@lucide/svelte';
	import type { DateRange } from 'bits-ui';
	import { RangeCalendar } from 'bits-ui';

	const df = new DateFormatter('en-US', {
		dateStyle: 'medium'
	});

	let {
		value = $bindable(),
		start,
		end,
		className
	} = $props<{
		value?: DateRange;
		start?: string;
		end?: string;
		className?: string;
	}>();

	// Initialize value from props if provided
	$effect(() => {
		if (!value && start) {
			value = {
				start: start
					? new CalendarDate(
							new Date(start).getFullYear(),
							new Date(start).getMonth() + 1,
							new Date(start).getDate()
						)
					: undefined,
				end: end
					? new CalendarDate(
							new Date(end).getFullYear(),
							new Date(end).getMonth() + 1,
							new Date(end).getDate()
						)
					: undefined
			};
		}
	});

	let open = $state(false);
	let placeholder = $state<CalendarDate | undefined>();

	// Initialize placeholder from start date
	$effect(() => {
		if (value?.start && !placeholder) {
			placeholder = value.start;
		}
	});
</script>

<div class={cn('grid gap-2', className)}>
	<Popover.Root bind:open>
		<Popover.Trigger>
			{#snippet child({ props })}
				<Button
					variant="outline"
					class={cn(
						'w-[300px] justify-start text-left font-normal',
						!value && 'text-muted-foreground'
					)}
					{...props}
				>
					<CalendarIcon class="mr-2 h-4 w-4" />
					{#if value && value.start}
						{#if value.end}
							{df.format(value.start.toDate(getLocalTimeZone()))} - {df.format(
								value.end.toDate(getLocalTimeZone())
							)}
						{:else}
							{df.format(value.start.toDate(getLocalTimeZone()))}
						{/if}
					{:else}
						Pick a date
					{/if}
				</Button>
			{/snippet}
		</Popover.Trigger>
		<Popover.Content class="w-auto p-0" align="start">
			<RangeCalendar.Root bind:value bind:placeholder class="rounded-md border p-3">
				{#snippet children({ months, weekdays })}
					<RangeCalendar.Header class="flex w-full items-center justify-between">
						<RangeCalendar.PrevButton
							class="h-9 w-9 bg-transparent p-0 opacity-50 hover:opacity-100"
						/>
						<RangeCalendar.Heading class="text-sm font-medium" />
						<RangeCalendar.NextButton
							class="h-9 w-9 bg-transparent p-0 opacity-50 hover:opacity-100"
						/>
					</RangeCalendar.Header>
					<div class="mt-4 flex flex-col gap-y-4 sm:flex-row sm:gap-x-4 sm:gap-y-0">
						{#each months as month}
							<RangeCalendar.Grid class="w-full border-collapse space-y-1">
								<RangeCalendar.GridHead>
									<RangeCalendar.GridRow class="mb-1 grid w-full grid-cols-7">
										{#each weekdays as weekday}
											<RangeCalendar.HeadCell
												class="w-9 rounded-md text-[0.8rem] font-normal text-muted-foreground"
											>
												{weekday.slice(0, 2)}
											</RangeCalendar.HeadCell>
										{/each}
									</RangeCalendar.GridRow>
								</RangeCalendar.GridHead>
								<RangeCalendar.GridBody>
									{#each month.weeks as weekDates}
										<RangeCalendar.GridRow class="flex w-full">
											{#each weekDates as date}
												<RangeCalendar.Cell
													{date}
													month={month.value}
													class="relative p-0 text-center text-sm focus-within:relative focus-within:z-20 [&:has([data-selected])]:bg-accent first:[&:has([data-selected])]:rounded-l-md last:[&:has([data-selected])]:rounded-r-md [&:has([data-selected][data-outside-month])]:bg-accent/50 [&:has([data-selected][data-selection-end])]:bg-accent [&:has([data-selected][data-selection-start])]:bg-accent"
												>
													<RangeCalendar.Day
														class={cn(
															'h-9 w-9 p-0 font-normal aria-selected:opacity-100',
															date.compare(month.value) !== 0 && 'text-muted-foreground opacity-50'
														)}
													/>
												</RangeCalendar.Cell>
											{/each}
										</RangeCalendar.GridRow>
									{/each}
								</RangeCalendar.GridBody>
							</RangeCalendar.Grid>
						{/each}
					</div>
				{/snippet}
			</RangeCalendar.Root>
		</Popover.Content>
	</Popover.Root>
</div>
